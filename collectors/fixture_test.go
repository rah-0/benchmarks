package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMakeRecords(t *testing.T) {
	for _, size := range []int{256, 1024} {
		records, lines, err := makeRecords(12, size)
		if err != nil {
			t.Fatal(err)
		}
		if len(records) != 13 || len(lines) != 13 {
			t.Fatalf("expected warmup and 12 records, got %d records and %d lines", len(records), len(lines))
		}
		for id, line := range lines {
			if len(line) != size || line[len(line)-1] != '\n' || bytes.Count(line, []byte{'\n'}) != 1 {
				t.Fatalf("record %d is not one %d-byte NDJSON line", id, size)
			}
			var decoded record
			if err := json.Unmarshal(line, &decoded); err != nil {
				t.Fatal(err)
			}
			if decoded.ID != id || decoded != records[id] {
				t.Fatalf("record %d did not round-trip", id)
			}
		}
	}
	if _, _, err := makeRecords(1, 1); err == nil {
		t.Fatal("accepted a record size too small for the required fields")
	}
}

func TestReceiverWarmupFormatsAndDuplicates(t *testing.T) {
	records, lines := fixtureRecords(t, 2)
	r := newReceiver(records)
	postFixture(t, r, string(lines[0]), http.StatusOK)
	stats, err := r.wait(t.Context(), true)
	if err != nil || stats.Unique != 0 || stats.Requests != 0 {
		t.Fatalf("warmup counted as measured work: %+v, %v", stats, err)
	}
	start := time.Now()
	r.begin(start, 0)
	array, err := json.Marshal(records[1:2])
	if err != nil {
		t.Fatal(err)
	}
	postFixture(t, r, " \n"+string(array), http.StatusOK)
	body := string(lines[2]) + string(lines[2]) + string(lines[1])
	response := postFixture(t, r, body, http.StatusOK)
	var ack struct {
		Status []struct{ Successful, Failed int }
	}
	if err := json.Unmarshal(response.Body.Bytes(), &ack); err != nil {
		t.Fatal(err)
	}
	if len(ack.Status) != 1 || ack.Status[0].Successful != 3 || ack.Status[0].Failed != 0 {
		t.Fatalf("wrong batch acknowledgement: %s", response.Body.String())
	}
	stats, err = r.wait(t.Context(), false)
	if err != nil || stats.Unique != 2 || stats.Duplicates != 2 || stats.Requests != 2 || stats.Finished.Before(start) {
		t.Fatalf("unexpected completion: %+v, %v", stats, err)
	}
}

func TestReceiverRejectsInvalidBatchesAtomically(t *testing.T) {
	records, lines := fixtureRecords(t, 2)
	first, second := string(lines[1]), string(lines[2])
	cases := map[string]string{
		"empty":          "",
		"empty array":    "[]",
		"malformed":      first + "{",
		"unknown field":  first + strings.Replace(second, `"id":2`, `"extra":true,"id":2`, 1),
		"nested field":   first + strings.Replace(second, `"service":`, `"extra":true,"service":`, 1),
		"altered value":  first + strings.Replace(second, `"INFO"`, `"ERROR"`, 1),
		"negative ID":    first + strings.Replace(second, `"id":2`, `"id":-1`, 1),
		"out of range":   first + strings.Replace(second, `"id":2`, `"id":3`, 1),
		"array trailing": "[" + first + "]{}",
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			r := newReceiver(records)
			r.begin(time.Now(), 0)
			postFixture(t, r, body, http.StatusBadRequest)
			stats, err := r.wait(t.Context(), false)
			if err == nil || stats.Unique != 0 || stats.Duplicates != 0 || !stats.Finished.IsZero() {
				t.Fatalf("invalid batch partially accepted: %+v, %v", stats, err)
			}
		})
	}
	t.Run("before timing", func(t *testing.T) {
		r := newReceiver(records)
		postFixture(t, r, first, http.StatusBadRequest)
	})
}

func TestReceiverOutageAndCancellation(t *testing.T) {
	records, lines := fixtureRecords(t, 2)
	r := newReceiver(records)
	r.begin(time.Now(), time.Hour)
	postFixture(t, r, string(lines[1]), http.StatusServiceUnavailable)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	stats, err := r.wait(ctx, false)
	if !errors.Is(err, context.Canceled) || stats.Unique != 0 || stats.Requests != 1 || stats.RetryRequests != 1 {
		t.Fatalf("outage accepted work or lost the cancellation: %+v, %v", stats, err)
	}
	r.begin(time.Now().Add(-time.Hour), 0)
	postFixture(t, r, string(lines[1]), http.StatusOK)
	stats, err = r.wait(ctx, false)
	if !errors.Is(err, context.Canceled) || !strings.Contains(err.Error(), "received 1/2 unique records") {
		t.Fatalf("cancellation did not report missing records: %+v, %v", stats, err)
	}
	postFixture(t, r, string(lines[2]), http.StatusOK)
	stats, err = r.wait(t.Context(), false)
	if err != nil || stats.Unique != 2 || stats.Requests != 3 || stats.RetryRequests != 1 || stats.Finished.IsZero() {
		t.Fatalf("receiver did not resume after outage: %+v, %v", stats, err)
	}
}

func fixtureRecords(t *testing.T, count int) ([]record, [][]byte) {
	t.Helper()
	records, lines, err := makeRecords(count, 256)
	if err != nil {
		t.Fatal(err)
	}
	return records, lines
}

func postFixture(t *testing.T, r *receiver, body string, status int) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/default/bench/_json", strings.NewReader(body)))
	if w.Code != status {
		t.Fatalf("HTTP status %d, want %d: %s", w.Code, status, w.Body.String())
	}
	return w
}
