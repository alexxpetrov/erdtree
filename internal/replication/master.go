package replication

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"connectrpc.com/connect"
	dbv1 "github.com/oleksiip-aiola/erdtree/gen/api/v1"
	"github.com/oleksiip-aiola/erdtree/gen/api/v1/dbv1connect"
	"github.com/oleksiip-aiola/erdtree/internal/wal"
)

type SlaveInfo struct {
	Address  string
	Client   dbv1connect.ErdtreeStoreClient
	LastSync int64
}

type MasterReplication struct {
	slaves       map[string]*SlaveInfo
	slavesMx     sync.RWMutex
	wal          *wal.WAL
	batchSize    int
	syncInterval time.Duration
}

var master *MasterReplication

func NewMasterReplication(wal *wal.WAL, syncInterval time.Duration, batchSize int) *MasterReplication {
	master = &MasterReplication{
		wal:          wal,
		batchSize:    batchSize,
		syncInterval: syncInterval,
		slaves:       make(map[string]*SlaveInfo),
	}

	go master.startSyncLoop()

	return master
}

func (m *MasterReplication) AddSlave(address string) error {
	master.slavesMx.Lock()
	defer master.slavesMx.Unlock()
	if _, exists := master.slaves[address]; exists {
		return fmt.Errorf("slave node exists: %s", address)
	}

	baseURL := "http://" + address // Use https:// for secure connections
	client := &http.Client{
		// You can customize the HTTP client here if needed
		// For example, add timeouts, custom transport, etc.
	}

	// Create the Connect-Go client
	erdtreeClient := dbv1connect.NewErdtreeStoreClient(
		client,
		baseURL,
		connect.WithGRPC(), // Use this if you want gRPC compatibility
	)

	master.slaves[address] = &SlaveInfo{
		Address:  address,
		Client:   erdtreeClient,
		LastSync: 0,
	}

	return nil
}

func (m *MasterReplication) RemoveSlave(address string) {
	master.slavesMx.Lock()
	defer master.slavesMx.Unlock()

	if _, exists := master.slaves[address]; exists {
		// if conn, ok := slave.Client.(*grpc.ClientConn); ok {
		// 	conn.Close()
		// }
		delete(master.slaves, address)
	}
}

func (m *MasterReplication) startSyncLoop() {
	ticker := time.NewTicker(master.syncInterval)

	defer ticker.Stop()

	for range ticker.C {

		master.syncAllSlaves()
	}
}

func (m *MasterReplication) syncAllSlaves() {
	master.slavesMx.RLock()
	defer master.slavesMx.RUnlock()

	for _, slave := range master.slaves {
		go master.syncSlave(slave)
	}
}

func (m *MasterReplication) syncSlave(slave *SlaveInfo) {
	entries, err := master.wal.GetEntriesSince(slave.LastSync)

	if err != nil {
		fmt.Printf("Error getting WAL entries for slave %s: %v\n", slave.Address, err)
		return
	}

	for i := 0; i < len(entries); i += master.batchSize {
		end := i + master.batchSize

		if end > len(entries) {
			end = len(entries)
		}

		batch := entries[i:end]

		if err := master.sendBatchToSlave(slave, batch); err != nil {
			fmt.Printf("Error syncing batch to slave %s: %v\n", slave.Address, err)
		}

		slave.LastSync = batch[len(batch)-1].Timestamp
	}
}

func (m *MasterReplication) sendBatchToSlave(slave *SlaveInfo, batch []*dbv1.LogEntry) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	request := connect.NewRequest(&dbv1.ReplicationRequest{
		Entries: make([]*dbv1.LogEntry, len(batch)),
	})

	for i, entry := range batch {
		request.Msg.Entries[i] = &dbv1.LogEntry{
			Timestamp: entry.Timestamp,
			Operation: entry.Operation,
			Key:       entry.Key,
			Value:     entry.Value,
			ExpiresAt: entry.ExpiresAt,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := slave.Client.Replicate(ctx, request)
	fmt.Println(err)
	if err != nil {
		// Handle different types of errors
		if connectErr, ok := err.(*connect.Error); ok {
			switch connectErr.Code() {
			case connect.CodeUnavailable:
				return fmt.Errorf("slave unavailable: %s", slave.Address)
			case connect.CodeDeadlineExceeded:
				return fmt.Errorf("replication timed out for slave: %s", slave.Address)
			default:
				return fmt.Errorf("replication failed for slave %s: %v", slave.Address, err)
			}
		}
		return fmt.Errorf("unknown error during replication to slave %s: %v", slave.Address, err)
	}

	// Check the response
	if !resp.Msg.Success {
		return fmt.Errorf("replication to slave %s failed: %s", slave.Address, resp.Msg.Error)
	}

	return nil
}

func (m *MasterReplication) ReplicateEntry(entry *dbv1.LogEntry) {
	master.slavesMx.RLock()
	defer master.slavesMx.RUnlock()

	for _, slave := range master.slaves {
		go func(slave *SlaveInfo) {
			if err := master.sendBatchToSlave(slave, []*dbv1.LogEntry{entry}); err != nil {
				fmt.Printf("Failed to replicate entry %s to %s", entry.Key, slave.Address)
			}
		}(slave)
	}
}
