package replication

import (
	"context"
	"fmt"
	"sync"
	"time"

	"connectrpc.com/connect"
	dbv1 "github.com/oleksiip-aiola/erdtree/gen/api/v1"
	"github.com/oleksiip-aiola/erdtree/internal/db"
	"github.com/oleksiip-aiola/erdtree/internal/wal"
)

type SlaveReplication struct {
	store       sync.Map
	wal         *wal.WAL
	lastApplied int64
	mx          sync.RWMutex
}

func NewSlaveReplication(wal *wal.WAL) *SlaveReplication {
	return &SlaveReplication{
		wal: wal,
	}
}

func (slave *SlaveReplication) Replicate(ctx context.Context, req *connect.Request[dbv1.ReplicationRequest]) (*connect.Response[dbv1.ReplicationResponse], error) {
	slave.mx.Lock()
	defer slave.mx.Unlock()
	for _, entry := range req.Msg.Entries {
		if entry.Timestamp <= slave.lastApplied {
			continue
		}
		if err := slave.applyEntry(entry); err != nil {
			res := connect.NewResponse(&dbv1.ReplicationResponse{
				Success: false,
				Error:   fmt.Sprintf("Failed to apply entry %v", err),
			})

			return res, nil
		}

		slave.lastApplied = entry.Timestamp
	}

	res := connect.NewResponse(&dbv1.ReplicationResponse{Success: true})

	return res, nil
}

func (slave *SlaveReplication) applyEntry(entry *dbv1.LogEntry) error {
	if err := slave.wal.AppendEntry(entry); err != nil {
		return fmt.Errorf("failed to append to WAL: %w", err)
	}

	switch entry.Operation {
	case dbv1.Operation_SET:

		slave.store.Store(entry.Key, &db.Value{
			Data:         entry.Value,
			LastAccessed: time.Now(),
			ExpiresAt:    time.Unix(0, entry.ExpiresAt),
		})

	case dbv1.Operation_DELETE:
		slave.store.Delete(entry.Key)
	default:
		return fmt.Errorf("Unknown operation %v", entry.Operation)
	}

	return nil
}

func (slave *SlaveReplication) GetLastAppliedTimestamp() int64 {
	slave.mx.RLock()
	defer slave.mx.RUnlock()
	return slave.lastApplied
}
