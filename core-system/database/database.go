package database

import "sync"

type MemoryDB struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		data: make(map[string]string),
	}
}

func (db *MemoryDB) Set(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[key] = value
}

func (db *MemoryDB) Get(key string) string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.data[key]
}