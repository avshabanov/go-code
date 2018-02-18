package main

import (
	"context"
	"log"
	"time"
)

// Value represents certain data being handled in a context-aware operation
type Value int64

func calcSlow(x int) int {
	time.Sleep(100 * time.Millisecond)
	return x + 1
}

// calls calcSlow and consults context on whether operation should be halted
func fetchSlowResults(c context.Context, countOfCalls int, resultSink chan<- Value) error {
	for k, i := 0, 0; i < countOfCalls; i++ {
		k = calcSlow(k) // iteration - call slow operation
		select {
		case <-c.Done(): // is it time to halt?
			return c.Err() // yes
		case resultSink <- Value(k):
			log.Printf("slowly fetched value=%d", k)
		}
	}
	return nil
}

func demoContextTimeout(c context.Context, resultSinkCapacity int, countOfCalls int, timeout time.Duration) {
	c, cancel := context.WithTimeout(c, timeout)
	defer cancel()
	log.Printf("demoContextTimeout: start with resultSinkCapacity=%d, countOfCalls=%d, timeout=%s", resultSinkCapacity, countOfCalls, timeout)

	resultSink := make(chan Value, resultSinkCapacity)
	if err := fetchSlowResults(c, countOfCalls, resultSink); err != nil {
		log.Printf("demoContextTimeout: completed with error: %v", err)
	} else {
		log.Println("demoContextTimeout: successfully completed")
	}
}

func demoContextManualCancellation(c context.Context) {
	c, cancel := context.WithCancel(c)
	log.Printf("demoContextManualCancellation: started")

	go func(c context.Context) {
		for {
			select {
			case <-c.Done():
				return
			case <-time.After(50 * time.Millisecond):
				log.Printf("still waiting for cancellation...")
			}
		}
	}(c)

	time.Sleep(200 * time.Millisecond)
	cancel()

	log.Printf("demoContextManualCancellation: completed")
}

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	root := context.Background()

	// 6*100ms > 400ms  => timeout
	demoContextTimeout(root, 9, 6, 400*time.Millisecond)

	// sink capacity is not enough to fit all results => timeout
	demoContextTimeout(root, 1, 2, 300*time.Millisecond)

	// 2*100ms < 300ms, sink capacity is enough to fit all results => successful completion
	demoContextTimeout(root, 2, 2, 300*time.Millisecond)

	demoContextManualCancellation(root)
}
