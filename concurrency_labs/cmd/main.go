package main

import (
	"concurrency_labs/patterns"
	"concurrency_labs/pipeline"
	"concurrency_labs/worker"
	"fmt"
	"time"
)

/*
main.go = Entry point of the program

Key Idea:
Go programs start execution from the main package → main() function.

We use this file to orchestrate demonstrations of concurrency patterns.
Each demo lives in its own package/file to illustrate Go packaging.
*/

func main() {
	fmt.Println("=======================================================")
	fmt.Println("  GO GOROUTINES & TRUE CONCURRENCY - COMPLETE GUIDE")
	fmt.Println("=======================================================")

	// SECTION 1 — Goroutines
	// Goroutines are lightweight threads managed by Go runtime.
	// They are NOT OS threads.
	// Very cheap to create (~2KB initial stack).
	fmt.Println("\n📌 SECTION 1: Basic Goroutines")
	fmt.Println("─────────────────────────────────────────────────────")
	patterns.DemoBasicGoroutines()

	// SECTION 2 — Channels
	// Channels enable communication between goroutines.
	// "Do not communicate by sharing memory; share memory by communicating."
	fmt.Println("\n📌 SECTION 2: Channels")
	fmt.Println("─────────────────────────────────────────────────────")
	patterns.DemoChannels()

	// SECTION 3 — Select
	// Select waits on multiple channel operations.
	fmt.Println("\n📌 SECTION 3: Select Statement")
	fmt.Println("─────────────────────────────────────────────────────")
	patterns.DemoSelect()

	// SECTION 4 — WaitGroups
	// WaitGroup synchronizes completion of goroutines.
	fmt.Println("\n📌 SECTION 4: WaitGroups")
	fmt.Println("─────────────────────────────────────────────────────")
	patterns.DemoWaitGroups()

	// SECTION 5 — Mutex
	// Protect shared memory from race conditions.
	fmt.Println("\n📌 SECTION 5: Mutex & Race Conditions")
	fmt.Println("─────────────────────────────────────────────────────")
	patterns.DemoMutex()

	// SECTION 6 — Worker Pool
	fmt.Println("\n📌 SECTION 6: Worker Pool Pattern")
	fmt.Println("─────────────────────────────────────────────────────")
	worker.DemoWorkerPool()

	// SECTION 7 — Pipeline
	fmt.Println("\n📌 SECTION 7: Pipeline Pattern")
	fmt.Println("─────────────────────────────────────────────────────")
	pipeline.DemoPipeline()

	// SECTION 8 — Fan-Out / Fan-In
	fmt.Println("\n📌 SECTION 8: Fan-Out / Fan-In Pattern")
	fmt.Println("─────────────────────────────────────────────────────")
	patterns.DemoFanOutFanIn()

	// SECTION 9 — Context
	fmt.Println("\n📌 SECTION 9: Context - Cancellation & Timeout")
	fmt.Println("─────────────────────────────────────────────────────")
	patterns.DemoContext()

	// SECTION 10 — Atomic
	fmt.Println("\n📌 SECTION 10: Atomic Operations")
	fmt.Println("─────────────────────────────────────────────────────")
	patterns.DemoAtomic()

	fmt.Println("\n✅ All concurrency demos complete!")
	fmt.Println("💡 TIP: Run with race detector:")
	fmt.Println("   go run -race ./cmd/main.go")

	// Small sleep ensures background goroutines finish printing
	time.Sleep(200 * time.Millisecond)
}
