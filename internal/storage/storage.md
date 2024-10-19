package storage

import (
	"errors"
	"fmt"
	"sync"

	dbv1 "github.com/oleksiip-aiola/erdtree/gen/api/v1"
	"github.com/oleksiip-aiola/erdtree/internal/wal"
)

type Storage struct {
	data map[string][]byte
	wal  *wal.WAL
	mx   sync.Mutex
}

func NewStorage(wal *wal.WAL) *Storage {
	return &Storage{
		data: make(map[string][]byte),
		wal:  wal,
	}
}

func (s *Storage) Set(key string, value []byte) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.data[key] = value

	return s.wal.AppendEntry([]byte(fmt.Sprintf("SET %s %s", key, value)))
}

func (s *Storage) Get(key string) []byte {
	s.mx.Lock()
	defer s.mx.Unlock()

	return s.data[key]
}

func (s *Storage) Delete(key string) {
	s.mx.Lock()
	defer s.mx.Unlock()

	delete(s.data, key)

	s.wal.AppendEntry([]byte(fmt.Sprintf("DELETE %s", key)))
}

func (s *Storage) Recover(key string) (bool, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	entries, err := s.wal.Recover()

	if err != nil {
		return false, errors.New("failed to recover WAL")
	}

	for _, entry := range entries {

		switch entry.Operation {
		case dbv1.Operation_SET:
			s.Set(entry.Key, entry.Value)
		case dbv1.Operation_DELETE:
			s.Delete(key)
		}
	}

	return true, nil
}
