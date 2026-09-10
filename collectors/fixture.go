package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type record struct {
	ID         int              `json:"id"`
	Level      string           `json:"level"`
	Message    string           `json:"message"`
	Attributes recordAttributes `json:"attributes"`
}

type recordAttributes struct {
	Service string  `json:"service"`
	Attempt int     `json:"attempt"`
	OK      bool    `json:"ok"`
	Value   float64 `json:"value"`
}

// Index zero is a warmup record. Sizes include the terminating newline.
func makeRecords(count, size int) ([]record, [][]byte, error) {
	records := make([]record, count+1)
	lines := make([][]byte, count+1)
	for id := range records {
		r := record{ID: id, Level: "INFO", Attributes: recordAttributes{
			Service: "catalog", Attempt: 2, OK: true, Value: 1.25,
		}}
		base, err := json.Marshal(r)
		if err != nil {
			return nil, nil, err
		}
		padding := size - len(base) - 1
		if padding < 0 {
			return nil, nil, fmt.Errorf("record size %d is too small for record %d", size, id)
		}
		r.Message = strings.Repeat("x", padding)
		encoded, err := json.Marshal(r)
		if err != nil {
			return nil, nil, err
		}
		records[id] = r
		lines[id] = append(encoded, '\n')
	}
	return records, lines, nil
}

type receiptStats struct {
	Unique        int
	Duplicates    int
	Requests      int
	RetryRequests int
	Finished      time.Time
}

type receiver struct {
	mu          sync.Mutex
	expected    []record
	seen        []bool
	warm        bool
	started     bool
	unavailable time.Time
	stats       receiptStats
	err         error
	changed     chan struct{}
}

func newReceiver(expected []record) *receiver {
	return &receiver{expected: expected, seen: make([]bool, len(expected)), changed: make(chan struct{}, 1)}
}

func (r *receiver) begin(start time.Time, outage time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.started = true
	r.unavailable = start.Add(outage)
}

func (r *receiver) notify() {
	select {
	case r.changed <- struct{}{}:
	default:
	}
}

func (r *receiver) fail(w http.ResponseWriter, err error) {
	r.mu.Lock()
	r.err = err
	r.mu.Unlock()
	http.Error(w, "invalid benchmark payload", http.StatusBadRequest)
	r.notify()
}

func (r *receiver) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost || req.URL.Path != "/api/default/bench/_json" {
		http.NotFound(w, req)
		return
	}
	r.mu.Lock()
	retry := r.started && time.Now().Before(r.unavailable)
	if r.started {
		r.stats.Requests++
		if retry {
			r.stats.RetryRequests++
		}
	}
	r.mu.Unlock()
	if retry {
		// Drain rejected bodies to keep HTTP connection reuse comparable.
		_, _ = io.Copy(io.Discard, http.MaxBytesReader(w, req.Body, 16<<20))
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	reader := bufio.NewReader(http.MaxBytesReader(w, req.Body, 16<<20))
	first, err := reader.ReadByte()
	for err == nil && strings.ContainsRune(" \t\r\n", rune(first)) {
		first, err = reader.ReadByte()
	}
	if err != nil {
		r.fail(w, fmt.Errorf("read request: %w", err))
		return
	}
	_ = reader.UnreadByte()
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	var records []record
	if first == '[' {
		err = decoder.Decode(&records)
		if err == nil {
			var trailing any
			if extra := decoder.Decode(&trailing); !errors.Is(extra, io.EOF) {
				err = errors.New("data after JSON array")
			}
		}
	} else {
		for {
			var record record
			err = decoder.Decode(&record)
			if errors.Is(err, io.EOF) {
				err = nil
				break
			}
			if err != nil {
				break
			}
			records = append(records, record)
		}
	}
	if err != nil || len(records) == 0 {
		r.fail(w, fmt.Errorf("decode nonempty request: %v", err))
		return
	}
	r.mu.Lock()
	for _, record := range records {
		if record.ID < 0 || record.ID >= len(r.expected) || record != r.expected[record.ID] {
			r.mu.Unlock()
			r.fail(w, fmt.Errorf("unexpected or altered record %d", record.ID))
			return
		}
		if record.ID != 0 && !r.started {
			r.mu.Unlock()
			r.fail(w, errors.New("measured record arrived before timing started"))
			return
		}
	}
	// Acknowledge only fully validated batches. The completion time includes
	// writing and flushing the HTTP acknowledgement, not collector checkpointing.
	w.Header().Set("Content-Type", "application/json")
	_, err = fmt.Fprintf(w, `{"code":200,"status":[{"successful":%d,"failed":0}]}`, len(records))
	if err == nil {
		err = http.NewResponseController(w).Flush()
	}
	if err != nil {
		r.err = fmt.Errorf("acknowledge request: %w", err)
	} else {
		for _, record := range records {
			if record.ID == 0 {
				r.warm = true
				continue
			}
			if r.seen[record.ID] {
				r.stats.Duplicates++
			} else {
				r.seen[record.ID] = true
				r.stats.Unique++
			}
		}
		if r.stats.Unique == len(r.expected)-1 && r.stats.Finished.IsZero() {
			r.stats.Finished = time.Now()
		}
	}
	r.mu.Unlock()
	r.notify()
}

func (r *receiver) wait(ctx context.Context, warmup bool) (receiptStats, error) {
	for {
		r.mu.Lock()
		stats, err := r.stats, r.err
		done := !stats.Finished.IsZero()
		if warmup {
			done = r.warm
		}
		r.mu.Unlock()
		if err != nil || done {
			return stats, err
		}
		select {
		case <-ctx.Done():
			return stats, fmt.Errorf("received %d/%d unique records: %w", stats.Unique, len(r.expected)-1, ctx.Err())
		case <-r.changed:
		}
	}
}
