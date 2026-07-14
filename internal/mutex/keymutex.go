package mutex

import (
	"sync"
)

// KeyMutex provides a way to lock based on a string key.
// This is useful for ensuring only one operation happens per resource identifier at a time.
type KeyMutex struct {
	locks map[string]*sync.Mutex
	mu    sync.Mutex
}

// NewKeyMutex creates a new KeyMutex instance.
func NewKeyMutex() *KeyMutex {
	return &KeyMutex{
		locks: make(map[string]*sync.Mutex),
	}
}

// Lock acquires a lock for the given key.
// If no lock exists for this key, one is created.
func (km *KeyMutex) Lock(key string) {
	km.mu.Lock()
	if km.locks[key] == nil {
		km.locks[key] = &sync.Mutex{}
	}
	lock := km.locks[key]
	km.mu.Unlock()

	lock.Lock()
}

// Unlock releases the lock for the given key.
func (km *KeyMutex) Unlock(key string) {
	km.mu.Lock()
	lock := km.locks[key]
	km.mu.Unlock()

	if lock != nil {
		lock.Unlock()
	}
}
