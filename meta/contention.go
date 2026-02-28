package meta

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// ---------------------------------------------------------------------------
// 1. Atomic
// ---------------------------------------------------------------------------

type AtomicCounter struct {
	val atomic.Int64
}

func (c *AtomicCounter) Increment() {
	c.val.Add(1)
}

func (c *AtomicCounter) Get() int64 {
	return c.val.Load()
}

// ---------------------------------------------------------------------------
// 2. Mutex
// ---------------------------------------------------------------------------

type MutexCounter struct {
	mu  sync.Mutex
	val int64
}

func (c *MutexCounter) Increment() {
	c.mu.Lock()
	c.val++
	c.mu.Unlock()
}

func (c *MutexCounter) Get() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.val
}

// ---------------------------------------------------------------------------
// 3. RWMutex
// ---------------------------------------------------------------------------

type RWMutexCounter struct {
	mu  sync.RWMutex
	val int64
}

func (c *RWMutexCounter) Increment() {
	c.mu.Lock()
	c.val++
	c.mu.Unlock()
}

func (c *RWMutexCounter) Get() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.val
}

// ---------------------------------------------------------------------------
// 4. Channel
// ---------------------------------------------------------------------------

type channelReq struct {
	inc   int64
	reply chan int64 // nil for fire-and-forget ops
}

type ChannelCounter struct {
	ch   chan channelReq
	done chan struct{}
}

func NewChannelCounter(buf int) *ChannelCounter {
	c := &ChannelCounter{
		ch:   make(chan channelReq, buf),
		done: make(chan struct{}),
	}
	go func() {
		var v int64
		for {
			select {
			case r := <-c.ch:
				if r.inc != 0 {
					v += r.inc
				}
				if r.reply != nil {
					r.reply <- v
				}
			case <-c.done:
				return
			}
		}
	}()
	return c
}

func (c *ChannelCounter) Close() { close(c.done) }

func (c *ChannelCounter) Increment() { c.ch <- channelReq{inc: 1} }

// Get reads the current value. Caller provides reply channel to avoid per-call allocation.
func (c *ChannelCounter) Get(reply chan int64) int64 {
	c.ch <- channelReq{reply: reply}
	return <-reply
}

// ---------------------------------------------------------------------------
// 5. CAS Spin Loop (no backoff) — cautionary baseline showing raw CAS
//    failure cost under contention vs atomic.Add
// ---------------------------------------------------------------------------

type CASSpinCounter struct {
	val atomic.Int64
}

func (c *CASSpinCounter) Increment() {
	for {
		old := c.val.Load()
		if c.val.CompareAndSwap(old, old+1) {
			return
		}
	}
}

func (c *CASSpinCounter) Get() int64 {
	return c.val.Load()
}

// ---------------------------------------------------------------------------
// 5b. CAS with exponential backoff — demonstrates how yielding between
//     CAS retries reduces cache-line thrashing under high contention
// ---------------------------------------------------------------------------

type CASBackoffCounter struct {
	val atomic.Int64
}

func (c *CASBackoffCounter) Increment() {
	spins := 1
	for {
		old := c.val.Load()
		if c.val.CompareAndSwap(old, old+1) {
			return
		}
		for j := 0; j < spins; j++ {
			runtime.Gosched()
		}
		if spins < 16 {
			spins <<= 1
		}
	}
}

func (c *CASBackoffCounter) Get() int64 {
	return c.val.Load()
}

// ---------------------------------------------------------------------------
// 6. Ring Buffer (SPSC — single producer, single consumer)
//
// Only one goroutine may call Push and only one may call Pop.
// Element is written before head is published via atomic store, ensuring
// the consumer always reads a fully written slot.
// ---------------------------------------------------------------------------

type SPSCRingBuffer struct {
	buf  []int64
	mask int64
	head atomic.Int64 // written by producer, read by consumer
	tail atomic.Int64 // written by consumer, read by producer
}

func NewSPSCRingBuffer(size int) *SPSCRingBuffer {
	// round up to power of two for mask trick
	n := 1
	for n < size {
		n <<= 1
	}
	return &SPSCRingBuffer{
		buf:  make([]int64, n),
		mask: int64(n - 1),
	}
}

// Push enqueues a value. Returns false if the buffer is full.
// Must be called from a single producer goroutine.
func (r *SPSCRingBuffer) Push(val int64) bool {
	head := r.head.Load()
	next := (head + 1) & r.mask
	if next == r.tail.Load() {
		return false // full
	}
	r.buf[head] = val
	r.head.Store(next) // publish: element written before head advances
	return true
}

// Pop dequeues a value. Returns (0, false) if the buffer is empty.
// Must be called from a single consumer goroutine.
func (r *SPSCRingBuffer) Pop() (int64, bool) {
	tail := r.tail.Load()
	if tail == r.head.Load() {
		return 0, false // empty
	}
	val := r.buf[tail]
	r.tail.Store((tail + 1) & r.mask) // publish: element read before tail advances
	return val, true
}

// ---------------------------------------------------------------------------
// 7. sync.Map
// ---------------------------------------------------------------------------

type SyncMapStore struct {
	m sync.Map
}

func (s *SyncMapStore) Store(key string, val int64) {
	s.m.Store(key, val)
}

func (s *SyncMapStore) Load(key string) (int64, bool) {
	v, ok := s.m.Load(key)
	if !ok {
		return 0, false
	}
	return v.(int64), true
}

// MutexMap is the baseline comparison for sync.Map
type MutexMap struct {
	mu sync.RWMutex
	m  map[string]int64
}

func NewMutexMap() *MutexMap {
	return &MutexMap{m: make(map[string]int64)}
}

func (s *MutexMap) Store(key string, val int64) {
	s.mu.Lock()
	s.m[key] = val
	s.mu.Unlock()
}

func (s *MutexMap) Load(key string) (int64, bool) {
	s.mu.RLock()
	v, ok := s.m[key]
	s.mu.RUnlock()
	return v, ok
}

// ---------------------------------------------------------------------------
// 8. Semaphore (channel-based)
// ---------------------------------------------------------------------------

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(maxConcurrency int) *Semaphore {
	return &Semaphore{ch: make(chan struct{}, maxConcurrency)}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}

// ---------------------------------------------------------------------------
// 9. WaitGroup — thin wrapper to expose a consistent Work pattern
// ---------------------------------------------------------------------------

func WaitGroupWork(workers int, work func()) {
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			work()
		}()
	}
	wg.Wait()
}

// ---------------------------------------------------------------------------
// 10. sync.Pool
// ---------------------------------------------------------------------------

var Pool = sync.Pool{
	New: func() any {
		buf := make([]byte, 1024)
		return &buf
	},
}

func PoolGetPut() {
	buf := Pool.Get().(*[]byte)
	(*buf)[0] = 1 // simulate use
	Pool.Put(buf)
}

// ---------------------------------------------------------------------------
// 11. sync.Once
// ---------------------------------------------------------------------------

type OnceInit struct {
	once sync.Once
	val  int64
}

func (o *OnceInit) Init() {
	o.once.Do(func() {
		o.val = 42
	})
}

func (o *OnceInit) Get() int64 {
	return o.val
}

// ===========================================================================
// Contention-Reduction Algorithms
// ===========================================================================

// ---------------------------------------------------------------------------
// 12. Sharded / Striped Counter
//
// Splits one hot variable into N independent shards. Each writer picks a
// shard (round-robin via goroutine-local index) and increments it with
// atomic.Add. Reads aggregate all shards — slower but contention-free writes.
// ---------------------------------------------------------------------------

type paddedAtomicInt64 struct {
	v atomic.Int64
	_ [56]byte // pad to 64-byte cache line
}

type ShardedCounter struct {
	shards []paddedAtomicInt64
}

func NewShardedCounter(n int) *ShardedCounter {
	if n <= 0 {
		n = runtime.NumCPU()
	}
	return &ShardedCounter{shards: make([]paddedAtomicInt64, n)}
}

func (c *ShardedCounter) Increment(shard int) {
	c.shards[shard%len(c.shards)].v.Add(1)
}

func (c *ShardedCounter) Get() int64 {
	var total int64
	for i := range c.shards {
		total += c.shards[i].v.Load()
	}
	return total
}

// ---------------------------------------------------------------------------
// 13. Per-CPU Counter
//
// One counter per GOMAXPROCS slot. Writers use a slot index (provided by
// caller, typically derived from worker ID) to avoid cross-core contention.
// Reads sum all slots.
// ---------------------------------------------------------------------------

type PerCPUCounter struct {
	slots []paddedAtomicInt64
}

func NewPerCPUCounter() *PerCPUCounter {
	n := runtime.GOMAXPROCS(0)
	return &PerCPUCounter{slots: make([]paddedAtomicInt64, n)}
}

func (c *PerCPUCounter) Increment(slot int) {
	c.slots[slot%len(c.slots)].v.Add(1)
}

func (c *PerCPUCounter) Get() int64 {
	var total int64
	for i := range c.slots {
		total += c.slots[i].v.Load()
	}
	return total
}

// ---------------------------------------------------------------------------
// 14. Flat Combining
//
// Threads enqueue their operation into a shared list. One thread becomes the
// combiner, applies all pending ops in a batch, then publishes results.
// Others spin-wait for their result.
// ---------------------------------------------------------------------------

type combineOp struct {
	delta  int64
	done   atomic.Bool
	result int64
}

var combineOpPool = sync.Pool{
	New: func() any { return &combineOp{} },
}

type FlatCombiningCounter struct {
	mu      sync.Mutex
	val     int64
	pending []*combineOp
	pmu     sync.Mutex
}

func (c *FlatCombiningCounter) Increment() {
	op := combineOpPool.Get().(*combineOp)
	op.delta = 1
	op.done.Store(false)
	op.result = 0

	c.pmu.Lock()
	c.pending = append(c.pending, op)
	c.pmu.Unlock()

	// Try to become the combiner
	if c.mu.TryLock() {
		c.combine()
		c.mu.Unlock()
	} else {
		// Wait for combiner to process our op
		for !op.done.Load() {
			runtime.Gosched()
		}
	}
	combineOpPool.Put(op)
}

func (c *FlatCombiningCounter) combine() {
	c.pmu.Lock()
	ops := c.pending
	c.pending = nil // new backing array to avoid race with iteration
	c.pmu.Unlock()

	for _, op := range ops {
		c.val += op.delta
		op.result = c.val
		op.done.Store(true)
	}
}

func (c *FlatCombiningCounter) Get() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.combine()
	return c.val
}

// ---------------------------------------------------------------------------
// 15. Fetch-and-Add with Local Buffering
//
// Each goroutine accumulates increments in a local buffer. When the buffer
// reaches a threshold, it flushes to the global atomic counter in one Add.
// Trades eventual consistency for much higher write throughput.
// ---------------------------------------------------------------------------

type LocalBufferedCounter struct {
	global    atomic.Int64
	threshold int64
}

func NewLocalBufferedCounter(threshold int64) *LocalBufferedCounter {
	if threshold <= 0 {
		threshold = 64
	}
	return &LocalBufferedCounter{threshold: threshold}
}

// IncrementBuffered increments a caller-owned local buffer. When the buffer
// reaches the threshold, it flushes to the global counter and resets.
// Returns the updated local buffer value.
func (c *LocalBufferedCounter) IncrementBuffered(local *int64) {
	(*local)++
	if *local >= c.threshold {
		c.global.Add(*local)
		*local = 0
	}
}

// Flush forces any remaining local buffer to the global counter.
func (c *LocalBufferedCounter) Flush(local *int64) {
	if *local > 0 {
		c.global.Add(*local)
		*local = 0
	}
}

func (c *LocalBufferedCounter) Get() int64 {
	return c.global.Load()
}

// ---------------------------------------------------------------------------
// 16. MPSC Queue + Aggregator
//
// Multiple producers push increments into a lock-free MPSC queue (channel).
// A single aggregator goroutine drains the queue and updates the counter.
// Writers never touch shared state directly.
// ---------------------------------------------------------------------------

type MPSCAggregator struct {
	ch   chan int64
	val  atomic.Int64
	done chan struct{}
}

func NewMPSCAggregator(buf int) *MPSCAggregator {
	if buf <= 0 {
		buf = 4096
	}
	a := &MPSCAggregator{
		ch:   make(chan int64, buf),
		done: make(chan struct{}),
	}
	go a.loop()
	return a
}

func (a *MPSCAggregator) loop() {
	var v int64
	for {
		select {
		case delta := <-a.ch:
			v += delta
			// Drain as many as available without blocking
			for {
				select {
				case d := <-a.ch:
					v += d
				default:
					goto publish
				}
			}
		publish:
			a.val.Store(v)
		case <-a.done:
			// Drain remaining
			for {
				select {
				case d := <-a.ch:
					v += d
				default:
					a.val.Store(v)
					return
				}
			}
		}
	}
}

func (a *MPSCAggregator) Close() { close(a.done) }

func (a *MPSCAggregator) Increment() bool {
	select {
	case a.ch <- 1:
		return true
	case <-a.done:
		return false
	}
}

func (a *MPSCAggregator) Get() int64 { return a.val.Load() }

// ---------------------------------------------------------------------------
// 17. RCU (Read-Copy-Update)
//
// Readers load a pointer atomically — no locks, no contention.
// Writers copy the current state, modify the copy, then atomically swap
// the pointer. Old versions are left for GC (no explicit grace period
// needed in Go thanks to the GC).
// ---------------------------------------------------------------------------

type rcuState struct {
	val int64
}

type RCUCounter struct {
	state atomic.Pointer[rcuState]
}

func NewRCUCounter() *RCUCounter {
	c := &RCUCounter{}
	c.state.Store(&rcuState{})
	return c
}

func (c *RCUCounter) Increment() {
	for {
		old := c.state.Load()
		next := &rcuState{val: old.val + 1}
		if c.state.CompareAndSwap(old, next) {
			return
		}
	}
}

func (c *RCUCounter) Get() int64 {
	return c.state.Load().val
}

// ---------------------------------------------------------------------------
// 18. Disruptor-style MPSC Ring Buffer
//
// Preallocated ring buffer with sequence-based coordination. Multiple
// producers claim slots via atomic CAS on the write cursor. A single
// consumer reads published slots. Cache-friendly, zero-allocation after init.
// ---------------------------------------------------------------------------

type paddedCursor struct {
	v atomic.Int64
	_ [56]byte // pad to separate cache line
}

type DisruptorMPSC struct {
	buf       []int64
	mask      int64
	writeCur  paddedCursor // next slot to claim (producers)
	committed paddedCursor // highest contiguously committed slot
	readCur   paddedCursor // next slot to read (consumer)
}

func NewDisruptorMPSC(size int) *DisruptorMPSC {
	n := 1
	for n < size {
		n <<= 1
	}
	return &DisruptorMPSC{
		buf:  make([]int64, n),
		mask: int64(n - 1),
	}
}

// Publish claims a slot and writes a value. Returns false if the buffer is full.
func (d *DisruptorMPSC) Publish(val int64) bool {
	for {
		wc := d.writeCur.v.Load()
		rc := d.readCur.v.Load()
		if wc-rc > d.mask {
			return false // full
		}
		if d.writeCur.v.CompareAndSwap(wc, wc+1) {
			d.buf[wc&d.mask] = val
			// Wait until we can advance committed cursor contiguously
			for !d.committed.v.CompareAndSwap(wc, wc+1) {
				runtime.Gosched()
			}
			return true
		}
	}
}

// Consume reads the next available value. Returns (0, false) if empty.
// Must be called from a single consumer goroutine.
func (d *DisruptorMPSC) Consume() (int64, bool) {
	rc := d.readCur.v.Load()
	if rc >= d.committed.v.Load() {
		return 0, false
	}
	val := d.buf[rc&d.mask]
	d.readCur.v.Store(rc + 1)
	return val, true
}

// ---------------------------------------------------------------------------
// 19. Approximate Counter (Probabilistic)
//
// Accepts small counting error in exchange for massive throughput.
// Each increment only touches the global atomic counter with probability
// 1/sampleRate. The Get() result is scaled by sampleRate.
// Uses a fast xorshift PRNG per call to avoid rand contention.
// ---------------------------------------------------------------------------

type ApproxCounter struct {
	val        atomic.Int64
	sampleRate int64
}

func NewApproxCounter(sampleRate int64) *ApproxCounter {
	if sampleRate <= 0 {
		sampleRate = 100
	}
	return &ApproxCounter{sampleRate: sampleRate}
}

// Increment probabilistically updates the counter. seed must be a non-zero
// caller-owned state that persists across calls (one per goroutine).
func (c *ApproxCounter) Increment(seed *uint64) {
	// xorshift64
	if *seed == 0 {
		*seed = 1
	}
	s := *seed
	s ^= s << 13
	s ^= s >> 7
	s ^= s << 17
	*seed = s
	if int64(s%uint64(c.sampleRate)) == 0 {
		c.val.Add(c.sampleRate)
	}
}

func (c *ApproxCounter) Get() int64 {
	return c.val.Load()
}