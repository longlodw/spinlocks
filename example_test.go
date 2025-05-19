package spinlocks_test

import (
	"fmt"
	"github.com/longlodw/spinlocks"
	"sync"
)

// ExampleSpinLock demonstrates how to use the SpinLock.
func ExampleSpinLock() {
	var lock spinlocks.SpinLock
	var wg sync.WaitGroup
	counter := 0

	// Launch multiple goroutines that increment counter under lock protection.
	for i := range 20 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Acquire the lock
			lock.Lock()
			defer lock.Unlock() // Ensure the lock is released

			// Critical section: modify shared resource
			counter++
			fmt.Printf("Goroutine %d incremented counter to %d\n", id, counter)
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Output is non-deterministic, but should include lines showing increments.
}

// ExampleRWSpinLock demonstrates how to use RWSpinLock in a concurrent context.
func ExampleRWSpinLock() {
	var rwLock spinlocks.RWSpinLock
	var wg sync.WaitGroup
	counter := 0

	// Launch multiple readers
	for i := range 20 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			rwLock.RLock()
			fmt.Printf("Reader %d sees counter: %d\n", id, counter)
			rwLock.RUnlock()
		}(i)
	}

	// Launch a writer
	wg.Add(1)
	go func() {
		defer wg.Done()

		rwLock.Lock()
		counter++
		fmt.Println("Writer incremented counter")
		rwLock.Unlock()
	}()

	// Wait for all goroutines to complete
	wg.Wait()

	// Output may vary, but should show readers seeing the counter before/after modification.
}
