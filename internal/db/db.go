package db

import (
	"errors"
	"sync"
	"time"

	dbv1 "github.com/oleksiip-aiola/erdtree/gen/api/v1"
	"github.com/oleksiip-aiola/erdtree/internal/wal"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyExpired  = errors.New("key has expired")
	WalNotCreated  = errors.New("failed to create WAL")
)

type Value struct {
	Data         []byte
	ExpiresAt    time.Time
	LastAccessed time.Time
}

type InMemoryDB struct {
	data       sync.Map
	wal        *wal.WAL
	maxSize    int
	size       int
	gcInterval time.Duration
	mx         sync.Mutex
}

type Config struct {
	MaxSize      int
	GCInterval   time.Duration
	SyncInterval time.Duration
	WAL          *wal.WAL
}

func NewInMemoryDb(config *Config, wal *wal.WAL) (*InMemoryDB, error) {
	store := &InMemoryDB{
		wal:        wal,
		maxSize:    config.MaxSize,
		gcInterval: config.GCInterval,
	}

	return store, nil
}

func (db *InMemoryDB) Get(key string) ([]byte, error) {
	val, ok := db.data.Load(key)

	if !ok {
		return nil, ErrKeyNotFound
	}

	value := val.(*Value)

	if !value.ExpiresAt.IsZero() && time.Now().After(value.ExpiresAt) {
		db.data.Delete(key)
		return nil, ErrKeyExpired
	}

	value.LastAccessed = time.Now()
	return value.Data, nil
}

func (db *InMemoryDB) Set(key string, value []byte, ttl time.Duration) error {
	expiresAt := time.Time{}
	if ttl > 0 {
		expiresAt.Add(ttl)
	}

	entry := &dbv1.LogEntry{
		Timestamp: time.Now().UnixNano(),
		Operation: dbv1.Operation_SET,
		Key:       key,
		Value:     value,
		ExpiresAt: expiresAt.UnixNano(),
	}
	if err := db.wal.AppendEntry(entry); err != nil {
		return err
	}

	db.data.Store(key, &Value{
		Data:         value,
		LastAccessed: time.Now(),
		ExpiresAt:    expiresAt,
	})
	db.mx.Lock()
	db.size++
	db.mx.Unlock()

	if db.size == db.maxSize {
		db.evictExpiredOrLast()
	}

	return nil
}

func (db *InMemoryDB) Delete(key string) error {
	_, ok := db.data.Load(key)

	if !ok {
		return ErrKeyNotFound
	}

	entry := &dbv1.LogEntry{
		Timestamp: time.Now().UnixNano(),
		Operation: dbv1.Operation_DELETE,
		Key:       key,
	}

	if err := db.wal.AppendEntry(entry); err != nil {
		return err
	}

	db.data.Delete(key)
	db.mx.Lock()
	db.size--
	db.mx.Unlock()

	return nil
}

func (db *InMemoryDB) Recover() error {
	entries, err := db.wal.Recover(0)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		switch entry.Operation {
		case dbv1.Operation_SET:
			db.data.Store(entry.Key, &Value{
				Data:         entry.Value,
				LastAccessed: time.Now(),
				ExpiresAt:    time.Unix(0, entry.ExpiresAt),
			})
		case dbv1.Operation_DELETE:
			db.data.Delete(entry.Key)
		}
	}

	return nil
}

func (db *InMemoryDB) evictExpiredOrLast() {
	var oldestKey interface{}
	var oldestAccess time.Time
	now := time.Now()

	db.data.Range(func(key, value interface{}) bool {
		val := value.(*Value)
		if oldestAccess.IsZero() || val.LastAccessed.Before(oldestAccess) {
			oldestAccess = val.LastAccessed
			oldestKey = key
		}
		if !val.ExpiresAt.IsZero() && now.After(val.ExpiresAt) {
			db.data.Delete(key)
			db.size--
		}
		return true
	})

	// Maybe this lock isn't necessary here. What are the odds that another goroutine enters
	// this piece between sync.Map's controlled .Range and .Delete? Leaving it here just in case
	db.mx.Lock()
	if db.maxSize == db.size {
		db.data.Delete(oldestKey)
	}
	db.mx.Unlock()
}

func (db *InMemoryDB) startGC() {
	ticker := time.NewTicker(db.gcInterval)
	defer ticker.Stop()

	for range ticker.C {
		db.runGC()
	}
}

func (db *InMemoryDB) runGC() {
	now := time.Now()

	db.data.Range(func(key, value interface{}) bool {
		val := value.(*Value)

		if !val.ExpiresAt.IsZero() && now.After(val.ExpiresAt) {
			db.data.Delete(key)
			db.size--
		}
		return true
	})
}

func (db *InMemoryDB) Close() {
	db.wal.Close()
}
