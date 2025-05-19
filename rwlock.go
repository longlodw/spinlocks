package spinlocks

import (
	"sync"
	"sync/atomic"
)

// RWSpinLock is a simple read-write spinlock implementation using atomic operations inspired Pererson's algorithm.
type RWSpinLock struct {
	wantsToWrite atomic.Bool
	readers      atomic.Int32
	writable     atomic.Bool
}

type rLocker struct {
	*RWSpinLock
}

type wLocker struct {
	*RWSpinLock
}

// TryRLock attempts to acquire the read lock without blocking.
// It returns true if the lock was acquired, false otherwise.
func (l *RWSpinLock) TryRLock() bool {
	l.readers.Add(1)
	l.writable.Store(true)
	if l.wantsToWrite.Load() && l.writable.Load() {
		l.readers.Add(-1)
		return false
	}
	return true
}

// TryLock attempts to acquire the write lock without blocking.
// It returns true if the lock was acquired, false otherwise.
func (l *RWSpinLock) TryLock() bool {
	if l.wantsToWrite.CompareAndSwap(false, true) {
		l.writable.Store(false)
		if l.readers.Load() > 0 && !l.writable.Load() {
			l.wantsToWrite.Store(false)
			return false
		}
		return true
	}
	return false
}

// RUnlock releases the read lock.
// Do not call this method if the lock was not acquired.
func (l *RWSpinLock) RUnlock() {
	l.readers.Add(-1)
}

// Unlock releases the write lock.
// It is a no-op if the lock is not held.
func (l *RWSpinLock) Unlock() {
	l.wantsToWrite.Store(false)
}

// Lock acquires the write lock, blocking until it is available.
func (l *RWSpinLock) Lock() {
	for l.wantsToWrite.CompareAndSwap(false, true) {
		// Spin-wait until the lock is acquired
	}
	l.writable.Store(false)
	for l.readers.Load() > 0 && !l.writable.Load() {
		// Spin-wait until the lock is acquired
	}
}

// RLock acquires the read lock, blocking until it is available.
func (l *RWSpinLock) RLock() {
	l.readers.Add(1)
	l.writable.Store(true)
	for l.wantsToWrite.Load() && l.writable.Load() {
		// Spin-wait until the lock is acquired
	}
}

// RLocker returns a read locker for the RWSpinLock.
func (l *RWSpinLock) RLocker() sync.Locker {
	return &rLocker{l}
}

// WLocker returns a write locker for the RWSpinLock.
func (l *RWSpinLock) WLocker() sync.Locker {
	return &wLocker{l}
}

func (l *rLocker) Lock() {
	l.RLock()
}

func (l *rLocker) Unlock() {
	l.RUnlock()
}

func (l *wLocker) Lock() {
	l.RWSpinLock.Lock()
}

func (l *wLocker) Unlock() {
	l.RWSpinLock.Unlock()
}
