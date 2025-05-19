package spinlocks

import (
	"sync/atomic"
)

// SpinLock is a simple spinlock implementation using atomic operations.
type SpinLock struct {
	locked atomic.Bool
}

// TryLock attempts to acquire the lock without blocking.
// It returns true if the lock was acquired, false otherwise.
func (l *SpinLock) TryLock() bool {
	return l.locked.CompareAndSwap(false, true)
}

// Unlock releases the lock.
// It is a no-op if the lock is not held.
func (l *SpinLock) Unlock() {
	l.locked.Store(false)
}

// Lock acquires the lock, blocking until it is available.
func (l *SpinLock) Lock() {
	for !l.TryLock() {
		// Spin-wait until the lock is acquired
	}
}
