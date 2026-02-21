# 🐹 The Complete Go Developer Guide

> Goroutines · Concurrency Patterns · API Development · Package Ecosystem

---

## 📋 Table of Contents

1. [Go Runtime & The Scheduler](#1-go-runtime--the-scheduler)
2. [Goroutines — Every Type Explained](#2-goroutines--every-type-explained)
3. [Channels — Goroutine Communication](#3-channels--goroutine-communication)
4. [Synchronization Primitives](#4-synchronization-primitives)
5. [Concurrency Patterns](#5-concurrency-patterns)
6. [API Development in Go](#6-api-development-in-go)
7. [Finding Go Packages (Go's npm/PyPI)](#7-finding-go-packages-gos-npmpypi)

---

## 1. Go Runtime & The Scheduler

Before understanding goroutines, you need to understand **what runs them**.

```
┌─────────────────────────────────────────────────────────────┐
│                      YOUR GO PROGRAM                        │
│                                                             │
│  Goroutine 1   Goroutine 2   Goroutine 3   Goroutine N...  │
│       │              │             │              │         │
│       └──────────────┴─────────────┴──────────────┘         │
│                            │                                │
│                   ┌────────▼────────┐                       │
│                   │  Go Scheduler   │  ← M:N Scheduler      │
│                   │  (goroutine     │    maps M goroutines   │
│                   │   queue mgmt)   │    onto N OS threads   │
│                   └────────┬────────┘                       │
│                            │                                │
│          ┌─────────────────┼─────────────────┐             │
│          │                 │                 │             │
│    ┌─────▼─────┐    ┌──────▼─────┐   ┌──────▼─────┐       │
│    │ OS Thread │    │ OS Thread  │   │ OS Thread  │       │
│    │   (M1)    │    │   (M2)     │   │   (M3)     │       │
│    └─────┬─────┘    └──────┬─────┘   └──────┬─────┘       │
│          │                 │                 │             │
│    ┌─────▼─────────────────▼─────────────────▼─────┐       │
│    │              CPU Cores (GOMAXPROCS)            │       │
│    └────────────────────────────────────────────────┘       │
└─────────────────────────────────────────────────────────────┘
```

### Key Facts

| Concept | Goroutine | OS Thread |
|---|---|---|
| Initial Stack | ~2 KB | ~2 MB |
| Stack Growth | Dynamic (grows/shrinks) | Fixed |
| Creation Cost | Microseconds | Milliseconds |
| Switch Cost | Very cheap (userspace) | Expensive (kernel) |
| Managed By | Go runtime | Operating system |
| Max Count | Millions possible | Thousands (OS limit) |

### GOMAXPROCS

```go
import "runtime"

// By default = number of CPU cores on your machine
// This is what enables TRUE PARALLELISM in Go
fmt.Println(runtime.GOMAXPROCS(0))  // read current value

runtime.GOMAXPROCS(4)  // set to 4 OS threads
```

> **True Parallelism**: With `GOMAXPROCS=4`, Go runs 4 goroutines **simultaneously** on 4 CPU cores at the exact same instant. This is different from concurrency (taking turns).

---

## 2. Goroutines — Every Type Explained

### 2.1 Fire-and-Forget Goroutine

The simplest type. Launch it and never wait for it.

```
main goroutine
     │
     ├──── go doWork() ──────────────────────────────► runs independently
     │                                                  (main doesn't wait)
     │
     ▼
continues immediately
```

```go
func main() {
    go doWork()          // launched — main doesn't wait
    go sendEmail()       // launched — main doesn't wait
    go logEvent()        // launched — main doesn't wait

    time.Sleep(1 * time.Second)  // crude: give goroutines time to run
    // ⚠️ Don't do this in production — use WaitGroup instead
}

func doWork() {
    fmt.Println("working in background")
}
```

**Use cases:** Background jobs, event logging, metrics collection, fire notifications.

---

### 2.2 Goroutine with WaitGroup (Synchronized)

You launch goroutines but need to **wait for all of them to finish** before proceeding.

```
main goroutine
     │
     │  wg.Add(3)
     │
     ├──── go task(1) ──── does work ──── wg.Done() ─┐
     ├──── go task(2) ──── does work ──── wg.Done() ─┤
     ├──── go task(3) ──── does work ──── wg.Done() ─┤
     │                                                │
     │  wg.Wait() ◄───────────────────────────────────┘
     │  (blocks here until all 3 call Done)
     │
     ▼
continues after all tasks done
```

```go
func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1)   // ← ALWAYS before go keyword, never inside goroutine
        i := i      // ← CRITICAL: capture loop variable (shadow it)
        go func() {
            defer wg.Done()   // ← decrements counter when goroutine exits
            fmt.Printf("Task %d done\n", i)
        }()
    }

    wg.Wait()  // ← blocks until counter reaches zero
    fmt.Println("All tasks complete")
}
```

**The loop variable trap:**
```go
// ❌ WRONG — all goroutines print the same final value of i
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)  // i is shared — prints 3,3,3
    }()
}

// ✅ CORRECT — each goroutine has its own copy
for i := 0; i < 3; i++ {
    i := i  // creates new variable scoped to this iteration
    go func() {
        fmt.Println(i)  // prints 0,1,2 (in any order)
    }()
}
```

**Use cases:** Parallel file processing, batch API calls, parallel DB queries, map-reduce.

---

### 2.3 Goroutine with Channel (Result-Returning)

Goroutine sends its result back through a channel.

```
main goroutine
     │
     │  ch := make(chan int)
     │
     ├──── go compute(ch) ──── does work ──── ch <- result ──┐
     │                                                        │
     │  result := <-ch  ◄──────────────────────────────────── ┘
     │  (blocks until goroutine sends)
     │
     ▼
uses result
```

```go
func main() {
    ch := make(chan int)

    go func() {
        total := 0
        for i := 1; i <= 100; i++ {
            total += i
        }
        ch <- total  // send result
    }()

    result := <-ch   // receive result (blocks until goroutine sends)
    fmt.Println("Sum:", result)  // 5050
}
```

**Multiple results with buffered channel:**
```go
func main() {
    results := make(chan string, 3)  // buffer = 3, no blocking

    go func() { results <- fetchFromDB() }()
    go func() { results <- fetchFromCache() }()
    go func() { results <- fetchFromAPI() }()

    for i := 0; i < 3; i++ {
        fmt.Println(<-results)
    }
}
```

**Use cases:** Parallel computations, concurrent API calls where you need results, async data fetching.

---

### 2.4 Long-Running Background Goroutine

A goroutine that runs **forever** (or until cancelled) doing periodic or reactive work.

```
main goroutine                    background goroutine
     │                                    │
     │  go startWorker(ctx)               │ for {
     │ ─────────────────────────────────► │     select {
     │                                    │     case <-ticker.C: doWork()
     │  continues running                 │     case <-ctx.Done(): return
     │                                    │     }
     │                              ...   │ }
     │
     │  cancel()  ────────────────────────► goroutine exits
     │
     ▼
```

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())

    go backgroundWorker(ctx)  // runs forever until we cancel

    time.Sleep(5 * time.Second)
    cancel()  // signal the goroutine to stop
    time.Sleep(100 * time.Millisecond)  // let it clean up
}

func backgroundWorker(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            fmt.Println("tick: doing periodic work")
        case <-ctx.Done():
            fmt.Println("worker: shutting down")
            return
        }
    }
}
```

**Use cases:** Health checkers, metrics collectors, cache refreshers, queue consumers, connection pool managers.

---

### 2.5 Goroutine with Context (Cancellable)

Goroutines that respect **cancellation signals** and **timeouts**.

```
WithCancel:                          WithTimeout:
                                     
context.Background()                 context.Background()
        │                                    │
   WithCancel()                        WithTimeout(5s)
        │                                    │
   ┌────┴────┐                         ┌────┴────┐
   │   ctx   │──── cancel() ──────►  stops   │   ctx   │──── 5s passes ──► stops
   │         │                          │         │
   └─────────┘                          └─────────┘
        │                                    │
   goroutine checks                    goroutine checks
   <-ctx.Done()                        <-ctx.Done()
```

```go
// ── WithCancel ──────────────────────────────────────────────
ctx, cancel := context.WithCancel(context.Background())
defer cancel()  // always defer, prevents goroutine leaks

go func(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return  // clean exit
        default:
            doWork()
        }
    }
}(ctx)

// ── WithTimeout ─────────────────────────────────────────────
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := fetchData(ctx)
// err == context.DeadlineExceeded if timeout hit

// ── WithDeadline ─────────────────────────────────────────────
deadline := time.Now().Add(5 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

// ── WithValue (pass data through call chain) ─────────────────
type key string
ctx = context.WithValue(ctx, key("requestID"), "req-abc-123")

// Retrieve anywhere in the call chain:
if id, ok := ctx.Value(key("requestID")).(string); ok {
    fmt.Println("Request:", id)
}
```

**Context Rules:**
- Always pass `ctx` as the **first parameter**
- Always `defer cancel()` — even for timeout contexts
- Never store context in a struct field
- Check `ctx.Done()` in loops and blocking operations

---

### 2.6 Goroutine Leak — What to Avoid

A goroutine leak happens when a goroutine is started but **never exits**.

```
❌ LEAK: goroutine blocks forever, nobody reads the channel
──────────────────────────────────────────────────────────
func leak() {
    ch := make(chan int)       // unbuffered

    go func() {
        ch <- expensiveWork()  // ← blocks FOREVER if nobody reads ch
    }()

    // function returns without reading from ch
    // goroutine is stuck, consuming memory forever
}


✅ FIX 1: use buffered channel
──────────────────────────────
ch := make(chan int, 1)  // goroutine can send without blocking


✅ FIX 2: use context for cancellation
──────────────────────────────────────
go func(ctx context.Context) {
    select {
    case ch <- expensiveWork():
    case <-ctx.Done():
        return  // exits cleanly
    }
}(ctx)


✅ FIX 3: always have a reader
──────────────────────────────
go func() { ch <- work() }()
result := <-ch  // make sure someone reads
```

---

### 2.7 Goroutine Summary Table

| Type | Keyword/Tool | Blocks? | Returns Value? | Use Case |
|---|---|---|---|---|
| Fire-and-forget | `go func()` | No | No | Background jobs, logging |
| Synchronized | `go` + `WaitGroup` | Yes (Wait) | No | Batch parallel work |
| Result-returning | `go` + `channel` | Yes (receive) | Yes | Parallel computations |
| Long-running | `go` + `select` loop | Never returns | No | Daemons, watchers |
| Cancellable | `go` + `context` | Until cancelled | Optional | HTTP handlers, timeouts |

---

## 3. Channels — Goroutine Communication

> **"Don't communicate by sharing memory; share memory by communicating."** — Go proverb

### 3.1 Channel Types

```
┌─────────────────────────────────────────────────────────────┐
│                    CHANNEL TYPES                            │
│                                                             │
│  Unbuffered          Buffered                Directional    │
│  make(chan T)        make(chan T, N)          chan<- T       │
│                                              <-chan T       │
│  ┌───┐               ┌───┬───┬───┐                         │
│  │   │               │ 1 │ 2 │ 3 │          send-only ──►  │
│  └───┘               └───┴───┴───┘          receive-only ◄─│
│                                                             │
│  Synchronous:         Async up to            Compile-time  │
│  sender blocks        capacity              type safety    │
│  until receiver                                             │
│  ready                                                      │
└─────────────────────────────────────────────────────────────┘
```

```go
// Unbuffered — synchronous handoff
ch := make(chan string)

// Buffered — async up to capacity
ch := make(chan string, 10)

// Directional — restrict usage
func send(out chan<- int)  { out <- 42 }   // can only send
func recv(in <-chan int)   { v := <-in }   // can only receive

// Check if closed
v, ok := <-ch
if !ok {
    fmt.Println("channel is closed and empty")
}

// Range over channel (exits when channel is closed)
for v := range ch {
    fmt.Println(v)
}

// Always close from the SENDER side, never receiver
close(ch)
```

### 3.2 Channel Operations Cheat Sheet

```
Operation          Unbuffered    Buffered(full)   Buffered(space)  Closed
─────────────────────────────────────────────────────────────────────────
Send (ch <- v)     Block         Block            OK               PANIC
Receive (<-ch)     Block         OK               Block            Zero+false
Close (close(ch))  OK            OK               OK               PANIC
```

### 3.3 Select Statement

```go
// select picks whichever case is ready.
// If multiple ready → picks one randomly.
// default → non-blocking (runs if nothing else ready).

select {
case msg := <-ch1:
    // ch1 had data
case ch2 <- value:
    // sent to ch2
case <-time.After(5 * time.Second):
    // timeout
case <-ctx.Done():
    // context cancelled
default:
    // nothing was ready — non-blocking
}
```

---

## 4. Synchronization Primitives

### 4.1 sync.Mutex vs sync.RWMutex

```
sync.Mutex                     sync.RWMutex
─────────────────              ──────────────────────────
Only one goroutine             Multiple readers OR one writer
at a time (read or write)

Lock()   → exclusive           Lock()   → exclusive write
Unlock() → release             Unlock() → release write
                               RLock()  → shared read
                               RUnlock()→ release read

Good for:                      Good for:
• Simple counters               • Maps read frequently
• Any shared mutation           • Config that rarely changes
                               • Caches
```

```go
type SafeMap struct {
    mu sync.RWMutex
    m  map[string]int
}

func (s *SafeMap) Set(k string, v int) {
    s.mu.Lock()           // exclusive — blocks all readers+writers
    defer s.mu.Unlock()
    s.m[k] = v
}

func (s *SafeMap) Get(k string) int {
    s.mu.RLock()          // shared — other readers run concurrently
    defer s.mu.RUnlock()
    return s.m[k]
}
```

### 4.2 sync.Once

```go
var once sync.Once
var instance *DB

// GetDB is safe to call from any number of goroutines.
// The function passed to Do runs exactly once, ever.
func GetDB() *DB {
    once.Do(func() {
        instance = connectToDB()  // runs only on first call
    })
    return instance
}
```

### 4.3 sync/atomic

```go
var counter int64
var flag    int32

// Atomic increment — no mutex needed for single values
atomic.AddInt64(&counter, 1)

// Atomic read
val := atomic.LoadInt64(&counter)

// Compare-And-Swap: only swaps if current == expected
// Building block for lock-free data structures
swapped := atomic.CompareAndSwapInt32(&flag, 0, 1)

// When to use atomic vs mutex:
// • atomic → single primitive (int, pointer), counters, flags (FASTER)
// • mutex  → structs, multiple fields, complex logic
```

### 4.4 sync.Map

```go
// sync.Map is safe for concurrent use without additional locking.
// Best when: many goroutines read, few write, keys are stable.
// For heavy writes, a regular map + RWMutex is often faster.

var m sync.Map

m.Store("key", "value")

val, ok := m.Load("key")
if ok {
    fmt.Println(val.(string))
}

m.Delete("key")

m.Range(func(k, v interface{}) bool {
    fmt.Println(k, v)
    return true  // return false to stop iteration
})
```

---

## 5. Concurrency Patterns

### 5.1 Worker Pool

> Limit concurrent goroutines to N workers consuming from a job queue.

```
                    ┌─────────────┐
                    │  Job Queue  │
                    │  (channel)  │
                    └──────┬──────┘
                           │
          ┌────────────────┼────────────────┐
          │                │                │
    ┌─────▼─────┐    ┌─────▼─────┐   ┌─────▼─────┐
    │ Worker 1  │    │ Worker 2  │   │ Worker 3  │
    │ goroutine │    │ goroutine │   │ goroutine │
    └─────┬─────┘    └─────┬─────┘   └─────┬─────┘
          │                │                │
          └────────────────┼────────────────┘
                           │
                    ┌──────▼──────┐
                    │   Results   │
                    │  (channel)  │
                    └─────────────┘
```

```go
func workerPool(numWorkers, numJobs int) {
    jobs    := make(chan int, numJobs)
    results := make(chan int, numJobs)

    // Launch workers
    var wg sync.WaitGroup
    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {      // blocks until job available
                results <- job * job     // process and send result
            }
        }()
    }

    // Send jobs
    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)  // signal: no more jobs → workers exit their range loop

    // Close results when all workers done
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect results
    for r := range results {
        fmt.Println(r)
    }
}
```

---

### 5.2 Pipeline

> Chain stages where each stage transforms data streaming through channels.

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ generate │───►│  square  │───►│  filter  │───►│  print   │
│  stage   │    │  stage   │    │  stage   │    │  stage   │
│          │    │          │    │          │    │          │
│ 1,2,3..  │    │ 1,4,9..  │    │ >20 only │    │ 25,36..  │
└──────────┘    └──────────┘    └──────────┘    └──────────┘
     goroutine       goroutine       goroutine
     (all run concurrently)
```

```go
func generate(ctx context.Context, nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            select {
            case out <- n:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

func square(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            select {
            case out <- n * n:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

// Compose the pipeline
ctx := context.Background()
nums    := generate(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
squared := square(ctx, nums)

for result := range squared {
    fmt.Println(result)
}
```

---

### 5.3 Fan-Out / Fan-In

> Distribute work to many goroutines (fan-out), merge results (fan-in).

```
              ┌──── Worker A ────┐
              │                 │
source ───────┼──── Worker B ────┼──── merged ──── consumer
              │                 │
              └──── Worker C ────┘

fan-out: one source → many workers
fan-in:  many results → one channel
```

```go
// Fan-Out: start multiple workers, each reads from same input channel
func fanOut(input <-chan int, numWorkers int) []<-chan int {
    channels := make([]<-chan int, numWorkers)
    for i := 0; i < numWorkers; i++ {
        channels[i] = worker(input)
    }
    return channels
}

// Fan-In: merge multiple channels into one
func fanIn(channels ...<-chan int) <-chan int {
    merged := make(chan int)
    var wg sync.WaitGroup

    wg.Add(len(channels))
    for _, ch := range channels {
        ch := ch
        go func() {
            defer wg.Done()
            for v := range ch {
                merged <- v
            }
        }()
    }

    go func() {
        wg.Wait()
        close(merged)
    }()

    return merged
}
```

---

### 5.4 Semaphore (Bounded Concurrency)

> Simpler than a worker pool. Limit max goroutines without a queue.

```
Semaphore capacity = 3

Time:   goroutine: G1  G2  G3  G4  G5
─────────────────────────────────────────
t=0:    Acquire   [G1][G2][G3]
t=1:    G4 wants to Acquire → BLOCKS (capacity full)
t=2:    G1 Releases → G4 Acquires [G2][G3][G4]
```

```go
type Semaphore chan struct{}

func NewSemaphore(n int) Semaphore   { return make(Semaphore, n) }
func (s Semaphore) Acquire()         { s <- struct{}{} }  // blocks when full
func (s Semaphore) Release()         { <-s }              // frees a slot

// Usage
sem := NewSemaphore(5)  // max 5 concurrent goroutines

for _, url := range urls {
    url := url
    go func() {
        sem.Acquire()
        defer sem.Release()
        fetch(url)  // max 5 of these run simultaneously
    }()
}
```

---

### 5.5 Pattern Comparison

```
Pattern          When to use                         Goroutine count
─────────────────────────────────────────────────────────────────────
Fire-forget      Don't need results, don't need       Unlimited
                 to wait

WaitGroup        Need to wait for all to finish,       Unlimited
                 don't need results back

Worker Pool      Large number of jobs, bounded          Fixed N
                 concurrency, results needed

Pipeline         Streaming transformation,              One per stage
                 composable stages

Fan-Out/Fan-In   Parallel execution of same            Configurable
                 work, merge results

Semaphore        Limit goroutines without              Unlimited but
                 a job queue                           bounded active
```

---

## 6. API Development in Go

### 6.1 Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        CLIENT                                   │
│                (Browser / Mobile / curl)                        │
└────────────────────────────┬────────────────────────────────────┘
                             │ HTTP Request
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                     MIDDLEWARE CHAIN                            │
│                                                                 │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────┐   │
│  │  Logger  │→ │   CORS   │→ │ Recover  │→ │  RequireAuth │   │
│  │          │  │          │  │ (panic)  │  │   (JWT)      │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────────┘   │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                         ROUTER                                  │
│                   (http.ServeMux)                               │
│                                                                 │
│  GET  /health              →  HealthHandler                     │
│  POST /api/v1/auth/login   →  AuthHandler.Login                 │
│  GET  /api/v1/products     →  ProductHandler.List               │
│  POST /api/v1/products     →  ProductHandler.Create             │
│  PUT  /api/v1/products/{id}→  ProductHandler.Update             │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                        HANDLERS                                 │
│                                                                 │
│  1. Parse & validate request body                               │
│  2. Call service/database layer                                 │
│  3. Write JSON response                                         │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    DATABASE LAYER                               │
│                                                                 │
│  • CRUD operations                                              │
│  • Thread-safe (sync.RWMutex or connection pool)                │
│  • Returns domain models, not raw SQL rows                      │
└─────────────────────────────────────────────────────────────────┘
```

---

### 6.2 Project Structure

```
myapi/
├── go.mod                          ← module definition + dependencies
├── go.sum                          ← dependency checksums (commit this!)
├── Dockerfile                      ← container build instructions
├── docker-compose.yml              ← multi-service orchestration
├── Makefile                        ← convenience commands
│
├── cmd/
│   └── server/
│       └── main.go                 ← entry point: wire everything together
│
└── internal/                       ← private (cannot be imported externally)
    ├── config/
    │   └── config.go               ← read env vars, no hardcoded secrets
    ├── models/
    │   └── models.go               ← pure data structs (User, Product, etc.)
    ├── database/
    │   └── database.go             ← CRUD operations, thread-safe
    ├── middleware/
    │   └── middleware.go           ← Logger, CORS, Auth, Recovery
    └── handlers/
        ├── auth.go                 ← Register, Login
        ├── products.go             ← CRUD handlers
        └── users.go                ← profile handlers
```

**Why `internal/`?** Go enforces that nothing outside your module can import packages under `internal/`. It's a compile-time protection of your implementation details.

---

### 6.3 The `go.mod` File Explained

```
module github.com/yourname/myapi    ← module path (unique identifier)

go 1.22                             ← minimum Go version

require (
    github.com/golang-jwt/jwt/v5 v5.2.1   ← dependency + version
    github.com/google/uuid v1.6.0
)
```

```bash
go mod init github.com/yourname/myapi   # create go.mod
go mod tidy                              # add missing, remove unused deps
go get github.com/some/package          # add a dependency
go get github.com/some/package@v1.2.3   # specific version
go get github.com/some/package@latest   # latest version
```

---

### 6.4 HTTP Routing (Go 1.22 stdlib)

```go
mux := http.NewServeMux()

// Go 1.22: METHOD /path syntax built into stdlib
mux.HandleFunc("GET /products",         handler.List)
mux.HandleFunc("POST /products",        handler.Create)
mux.HandleFunc("GET /products/{id}",    handler.Get)    // {id} = path variable
mux.HandleFunc("PUT /products/{id}",    handler.Update)
mux.HandleFunc("DELETE /products/{id}", handler.Delete)

// Read path variable in handler:
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")   // Go 1.22
    // ...
}
```

---

### 6.5 Middleware Pattern

```go
// Middleware signature — takes a handler, returns a handler
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares in order
func Chain(h http.Handler, m ...func(http.Handler) http.Handler) http.Handler {
    for i := len(m) - 1; i >= 0; i-- {
        h = m[i](h)
    }
    return h
}

// Example middleware
func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)                              // call next
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
    })
}

// Wire it up
handler := Chain(mux, Logger, CORS, Recover)
```

---

### 6.6 Request → Response Flow

```
POST /api/v1/auth/login
Content-Type: application/json
{"email":"alice@example.com","password":"secret"}

                    │
                    ▼ Logger middleware: start timer
                    ▼ CORS middleware: add headers
                    ▼ Recover middleware: wrap in defer/recover
                    ▼ Router: match route → Login handler
                    │
    ┌───────────────▼───────────────────────────────────────┐
    │ Login Handler                                         │
    │                                                       │
    │  1. json.NewDecoder(r.Body).Decode(&req)              │
    │     → {"email":"alice@...", "password":"secret"}      │
    │                                                       │
    │  2. Validate fields (empty check, format check)       │
    │                                                       │
    │  3. db.GetUserByEmail(req.Email)                      │
    │     → User{ID:"abc", Email:"alice@...", ...}          │
    │                                                       │
    │  4. bcrypt.CompareHashAndPassword(hash, password)     │
    │     → matches ✓                                       │
    │                                                       │
    │  5. GenerateJWT(user.ID, user.Email, secret)          │
    │     → "eyJhbGciOiJIUzI1NiJ9..."                       │
    │                                                       │
    │  6. JSON(w, 200, APIResponse{Success:true, Data:...}) │
    └───────────────────────────────────────────────────────┘
                    │
                    ▼ Logger middleware: log duration
                    │
HTTP 200 OK
{"success":true,"data":{"token":"eyJ...","user":{...}}}
```

---

### 6.7 JWT Authentication Flow

```
REGISTRATION / LOGIN:

Client                        Server
  │                              │
  │── POST /auth/login ─────────►│
  │   {email, password}          │
  │                              │ 1. Verify password (bcrypt)
  │                              │ 2. Create JWT payload:
  │                              │    {user_id, email, role, exp}
  │                              │ 3. Sign with HMAC-SHA256(payload, SECRET)
  │◄─ 200 {token: "eyJ..."} ─────│
  │                              │
  │ (stores token)               │


AUTHENTICATED REQUEST:

Client                        Server (RequireAuth middleware)
  │                              │
  │── GET /api/v1/products ─────►│
  │   Authorization: Bearer eyJ  │
  │                              │ 1. Extract token from header
  │                              │ 2. Parse JWT
  │                              │ 3. Verify signature with SECRET
  │                              │ 4. Check exp (not expired)
  │                              │ 5. Store claims in context
  │                              │ 6. Call next handler
  │◄─ 200 {products: [...]} ─────│
```

```go
// JWT structure
// Header.Payload.Signature
// eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiMTIzIn0.abc123

// Header (base64):  {"alg":"HS256","typ":"JWT"}
// Payload (base64): {"user_id":"123","email":"a@b.com","exp":1234567890}
// Signature:        HMAC-SHA256(header+"."+payload, secret)
```

---

### 6.8 Standard JSON Response Format

```go
// Always return a consistent envelope:
// Success: {"success":true, "data":{...}, "meta":{...}}
// Error:   {"success":false, "error":"message"}

type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}

// JSON tags:
// `json:"field"`            → serialize as "field"
// `json:"field,omitempty"`  → omit if zero/nil
// `json:"-"`                → NEVER serialize (use for passwords!)
```

---

### 6.9 HTTP Status Codes Reference

```
2xx — Success
────────────────────────────────────────
200 OK           → GET, PUT (resource updated)
201 Created      → POST (new resource created)
204 No Content   → DELETE (success, no body)

4xx — Client Error
────────────────────────────────────────
400 Bad Request      → invalid JSON, missing fields
401 Unauthorized     → missing/invalid JWT token
403 Forbidden        → authenticated but no permission
404 Not Found        → resource doesn't exist
409 Conflict         → duplicate (email already registered)
422 Unprocessable    → valid JSON but fails business rules

5xx — Server Error
────────────────────────────────────────
500 Internal Server Error → unexpected server failure
503 Service Unavailable   → DB down, dependency unreachable
```

---

### 6.10 Docker — Multi-Stage Build

```
Without multi-stage:          With multi-stage:
──────────────────            ─────────────────
Your image includes:          Your image has only:

• Go compiler (~400MB)        • Your binary (~8MB)
• Build tools                 • CA certificates
• Source code                 • Nothing else
• Module cache
• Test files

Total: ~800MB                 Total: ~15MB
```

```dockerfile
# Stage 1: Build (large, ~1GB)
FROM golang:1.22-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download          # ← cached layer: only re-runs if go.mod changes
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server/main.go
#   ↑ no C libs    ↑ target   ↑ strip debug symbols (shrinks binary ~30%)

# Stage 2: Final (tiny, ~15MB)
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /server   # only copy the binary
USER nonroot:nonroot                       # never run as root
ENTRYPOINT ["/server"]
```

```bash
# Build
docker build -t myapi .

# Run
docker run -p 8080:8080 -e JWT_SECRET=mysecret myapi

# Check image size
docker images myapi
```

---

### 6.11 Environment Variables Best Practice

```go
// config/config.go — read ALL config at startup, never at call time

type Config struct {
    Port      string
    JWTSecret string
    DBUrl     string
}

func Load() *Config {
    return &Config{
        Port:      getEnv("PORT", "8080"),          // default: 8080
        JWTSecret: getEnv("JWT_SECRET", ""),        // no default: required!
        DBUrl:     getEnv("DATABASE_URL", ""),
    }
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}
```

```bash
# Inject via Docker
docker run -e JWT_SECRET=supersecret -e PORT=9000 myapi

# Inject via docker-compose.yml
environment:
  - JWT_SECRET=${JWT_SECRET}   # reads from .env file
  - PORT=8080

# .env file (never commit to git!)
JWT_SECRET=supersecretkey
DATABASE_URL=postgres://user:pass@localhost/db
```

---

## 7. Finding Go Packages (Go's npm/PyPI)

### 7.1 The Official Package Registry

> 🌐 **https://pkg.go.dev**

This is Go's official package documentation and discovery site. Every public Go module is indexed here automatically.

```
pkg.go.dev
│
├── Search bar at top → type any keyword
│   e.g. "http router", "jwt", "postgres driver"
│
├── Each package page shows:
│   ├── README / documentation
│   ├── All exported functions, types, methods
│   ├── Source code links
│   ├── Versions (all tags)
│   ├── Dependencies
│   ├── Dependents (who uses it)
│   └── License
│
└── URL pattern:
    pkg.go.dev/github.com/some/package
    pkg.go.dev/github.com/some/package@v1.2.3  ← specific version
```

---

### 7.2 How Go Modules Work

```
Python:  package name ≠ import path   (pip install requests → import requests)
Node:    package name ≠ import path   (npm install axios → require('axios'))

Go:      package path IS the import   (go get github.com/gin-gonic/gin)
         The module path is the full  → import "github.com/gin-gonic/gin"
         GitHub/URL path              (no separate registry needed!)
```

```bash
# Install a package
go get github.com/gin-gonic/gin

# Install specific version
go get github.com/gin-gonic/gin@v1.9.1

# Install latest version
go get github.com/gin-gonic/gin@latest

# View all available versions
go list -m -versions github.com/gin-gonic/gin

# Remove unused dependencies
go mod tidy

# Show dependency graph
go mod graph

# Upgrade all dependencies
go get -u ./...

# Upgrade a specific package
go get -u github.com/gin-gonic/gin
```

---

### 7.3 Essential Packages by Category

#### 🌐 HTTP Routers / Frameworks

| Package | Stars | When to use |
|---|---|---|
| `net/http` (stdlib) | built-in | Simple APIs, full control, Go 1.22+ has method routing |
| `github.com/gin-gonic/gin` | ⭐ 78k | Fast, popular, middleware ecosystem |
| `github.com/go-chi/chi` | ⭐ 18k | Lightweight, stdlib-compatible, composable |
| `github.com/labstack/echo` | ⭐ 30k | Fast, clean API |
| `github.com/gofiber/fiber` | ⭐ 33k | Express.js-like, very fast |

```bash
go get github.com/gin-gonic/gin
go get github.com/go-chi/chi/v5
```

#### 🗄️ Databases

| Package | DB | Notes |
|---|---|---|
| `github.com/jackc/pgx/v5` | PostgreSQL | Fastest, native Go |
| `database/sql` + `github.com/lib/pq` | PostgreSQL | stdlib driver |
| `go.mongodb.org/mongo-driver` | MongoDB | Official driver |
| `github.com/redis/go-redis/v9` | Redis | Official client |
| `gorm.io/gorm` | Any SQL | Full ORM, ActiveRecord-style |
| `github.com/jmoiron/sqlx` | Any SQL | Thin wrapper, struct scanning |
| `github.com/pressly/goose` | Any SQL | DB migrations |

```bash
go get github.com/jackc/pgx/v5
go get gorm.io/gorm
go get gorm.io/driver/postgres
```

#### 🔐 Auth & Security

| Package | Use |
|---|---|
| `github.com/golang-jwt/jwt/v5` | JWT tokens |
| `golang.org/x/crypto/bcrypt` | Password hashing |
| `github.com/alexedwards/argon2id` | Modern password hashing |

```bash
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto
```

#### ✅ Validation

| Package | Use |
|---|---|
| `github.com/go-playground/validator/v10` | Struct tag validation |
| `github.com/asaskevich/govalidator` | String validation |

```go
// go-playground/validator example:
type User struct {
    Name  string `validate:"required,min=2,max=100"`
    Email string `validate:"required,email"`
    Age   int    `validate:"gte=0,lte=130"`
}

validate := validator.New()
err := validate.Struct(user)
```

#### 📝 Logging

| Package | Style |
|---|---|
| `log/slog` (stdlib Go 1.21+) | Structured, built-in |
| `github.com/rs/zerolog` | Zero-allocation, JSON |
| `go.uber.org/zap` | Very fast, structured |
| `github.com/sirupsen/logrus` | Feature-rich |

```go
// slog (built-in since Go 1.21) — prefer this for new projects
slog.Info("server started", "port", 8080, "env", "production")
slog.Error("request failed", "error", err, "path", r.URL.Path)
```

#### 🔧 Configuration

| Package | Use |
|---|---|
| `github.com/spf13/viper` | YAML/JSON/TOML/env config |
| `github.com/joho/godotenv` | Load `.env` files |
| `github.com/kelseyhightower/envconfig` | Struct tags for env vars |

#### 🧪 Testing

| Package | Use |
|---|---|
| `testing` (stdlib) | Unit tests, benchmarks |
| `net/http/httptest` (stdlib) | HTTP handler testing |
| `github.com/stretchr/testify` | Assertions, mocking |
| `github.com/vektra/mockery` | Auto-generate mocks |

```go
// net/http/httptest — test handlers without a real server
func TestProductList(t *testing.T) {
    req  := httptest.NewRequest("GET", "/products", nil)
    w    := httptest.NewRecorder()

    handler.List(w, req)

    resp := w.Result()
    assert.Equal(t, 200, resp.StatusCode)
}
```

#### 📊 Observability

| Package | Use |
|---|---|
| `github.com/prometheus/client_golang` | Prometheus metrics |
| `go.opentelemetry.io/otel` | OpenTelemetry tracing |

#### 🛠️ Utilities

| Package | Use |
|---|---|
| `github.com/google/uuid` | UUID generation |
| `github.com/spf13/cobra` | CLI apps |
| `github.com/gorilla/websocket` | WebSockets |
| `github.com/robfig/cron/v3` | Cron job scheduling |

---

### 7.4 Versioning — How Go Does It

```
Semantic Versioning: v MAJOR . MINOR . PATCH
                       │        │        │
                       │        │        └── Bug fixes (backward compatible)
                       │        └────────── New features (backward compatible)
                       └─────────────────── Breaking changes

v1.9.1   → stable release
v2.0.0   → breaking changes → import path CHANGES to /v2
v0.x.x   → unstable, no compatibility guarantees
```

**Major version in import path:**
```go
// v1 (original)
import "github.com/gin-gonic/gin"

// v2+ gets a new import path
import "github.com/some/pkg/v2"
import "github.com/some/pkg/v3"

// This means: go get github.com/some/pkg/v2
```

---

### 7.5 Reading the `go.sum` File

```
go.sum stores cryptographic hashes of every dependency.
Commit it to git. It ensures builds are reproducible and tamper-proof.

github.com/gin-gonic/gin v1.9.1 h1:4idEAncQnU5cB7BeOkPtxjfCSye0TRs...
github.com/gin-gonic/gin v1.9.1/go.mod h1:hPzmOoM197QpJMFJi...

Format: module version h1:hash
        module version/go.mod h1:hash

Never edit go.sum manually. Run go mod tidy to regenerate.
```

---

### 7.6 Useful `go` Commands Cheatsheet

```bash
# Module management
go mod init github.com/you/project    # initialize new module
go mod tidy                            # sync deps with code
go mod download                        # download all deps to cache
go mod verify                          # verify checksums
go mod why github.com/some/pkg         # why is this dep needed?

# Dependency management
go get github.com/pkg/name             # add/update dep
go get github.com/pkg/name@v1.2.3     # specific version
go get github.com/pkg/name@none        # remove dep

# Build & run
go run ./cmd/server/main.go            # run without building
go build -o bin/server ./cmd/server/   # compile binary
go install ./...                       # install binaries to $GOPATH/bin

# Test
go test ./...                          # run all tests
go test -v ./...                       # verbose
go test -race ./...                    # race detector
go test -cover ./...                   # coverage
go test -bench=. ./...                 # benchmarks

# Code quality
go vet ./...                           # static analysis
go fmt ./...                           # format code
goimports -w .                         # format + fix imports

# Explore packages
go doc github.com/some/package         # docs in terminal
go list -m all                         # list all dependencies
go list -m -versions github.com/gin-gonic/gin  # available versions
```

---

## 📦 Dependency Cache Location

```bash
# Go stores downloaded modules here (not in your project!)
$GOPATH/pkg/mod/

# Default GOPATH locations:
# Linux/macOS: ~/go
# Windows:     %USERPROFILE%\go

# So your module cache is at:
# ~/go/pkg/mod/github.com/gin-gonic/gin@v1.9.1/

# Clear cache:
go clean -modcache
```

---

## 🚀 Quick Start Commands Summary

```bash
# ── Project 1: Goroutines ─────────────────────────────────────
mkdir mygoroutines && cd mygoroutines
go mod init github.com/yourname/goroutines
# create your .go files
go run ./cmd/main.go
go run -race ./cmd/main.go    # with race detector

# ── Project 2: API ────────────────────────────────────────────
mkdir myapi && cd myapi
go mod init github.com/yourname/myapi
go get github.com/golang-jwt/jwt/v5
go get github.com/google/uuid
# create your .go files
go mod tidy
go run ./cmd/server/main.go

# ── Docker ───────────────────────────────────────────────────
docker build -t myapi .
docker run -p 8080:8080 -e JWT_SECRET=secret myapi
docker compose up --build

# ── Find packages ────────────────────────────────────────────
# Browser:  https://pkg.go.dev
# Terminal: go list -m -versions github.com/some/package
```

---

*Built with Go 1.22 · pkg.go.dev for packages · docker for deployment*