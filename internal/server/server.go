package server

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	dbv1 "github.com/oleksiip-aiola/erdtree/gen/api/v1"
	"github.com/oleksiip-aiola/erdtree/internal/config"
	"github.com/oleksiip-aiola/erdtree/internal/db"
	"github.com/oleksiip-aiola/erdtree/internal/replication"
	"github.com/oleksiip-aiola/erdtree/internal/wal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type KVStoreServer struct {
	// dbv1.UnimplementedKVStoreServer
	db         *db.InMemoryDB
	wal        *wal.WAL
	masterRepl *replication.MasterReplication
	slaveRepl  *replication.SlaveReplication
	isMaster   bool
}

type Value struct {
	Data      []byte
	ExpiresAt time.Time
}

func NewKVStoreServer(walDir string, isMaster bool, slaveAddresses []string, db *db.InMemoryDB, wal *wal.WAL) (*KVStoreServer, error) {
	server := &KVStoreServer{
		db:       db,
		wal:      wal,
		isMaster: isMaster,
	}
	cfg, _ := config.LoadConfig("config.yaml")

	if isMaster {

		server.masterRepl = replication.NewMasterReplication(wal, cfg.Master.SyncInterval, cfg.Master.BatchSize)
		for _, addr := range slaveAddresses {
			if err := server.masterRepl.AddSlave(addr); err != nil {
				return nil, fmt.Errorf("failed to add slave %s: %w", addr, err)
			}
		}
	} else {

		server.slaveRepl = replication.NewSlaveReplication(wal)
	}

	if err := server.db.Recover(); err != nil {
		return nil, fmt.Errorf("failed to recover data: %w", err)
	}

	return server, nil
}

func (server *KVStoreServer) Get(ctx context.Context, req *connect.Request[dbv1.GetRequest]) (*connect.Response[dbv1.GetResponse], error) {
	value, err := server.db.Get(req.Msg.Key)
	if err == db.ErrKeyNotFound {
		return nil, status.Error(codes.NotFound, "key not found")
	} else if err == db.ErrKeyExpired {
		return nil, status.Error(codes.NotFound, "key expired")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	res := connect.NewResponse(&dbv1.GetResponse{Value: value})

	return res, nil
}

func (server *KVStoreServer) Set(ctx context.Context, req *connect.Request[dbv1.SetRequest]) (*connect.Response[dbv1.SetResponse], error) {
	if !server.isMaster {
		return nil, status.Error(codes.PermissionDenied, "cannot write to slave")
	}
	err := server.db.Set(req.Msg.Key, req.Msg.Value, time.Duration(1000)*time.Second)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to set value")
	}

	if server.masterRepl != nil {
		server.masterRepl.ReplicateEntry(&dbv1.LogEntry{
			Timestamp: time.Now().UnixNano(),
			Operation: dbv1.Operation_SET,
			Key:       req.Msg.Key,
			Value:     req.Msg.Value,
			ExpiresAt: time.Now().Add(time.Duration(1000) * time.Second).UnixNano(),
		})
	}

	res := connect.NewResponse(&dbv1.SetResponse{})

	return res, nil
}

func (server *KVStoreServer) Delete(ctx context.Context, req *connect.Request[dbv1.DeleteRequest]) (*connect.Response[dbv1.DeleteResponse], error) {
	if !server.isMaster {
		return nil, status.Error(codes.PermissionDenied, "cannot delete from slave")
	}

	err := server.db.Delete(req.Msg.Key)
	if err == db.ErrKeyNotFound {
		return nil, status.Error(codes.NotFound, "key not found")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete key")
	}

	if server.masterRepl != nil {
		server.masterRepl.ReplicateEntry(&dbv1.LogEntry{
			Timestamp: time.Now().UnixNano(),
			Operation: dbv1.Operation_DELETE,
			Key:       req.Msg.Key,
		})
	}
	res := connect.NewResponse(&dbv1.DeleteResponse{})

	return res, nil
}

func (server *KVStoreServer) Replicate(ctx context.Context, req *connect.Request[dbv1.ReplicationRequest]) (*connect.Response[dbv1.ReplicationResponse], error) {
	if server.isMaster {
		return nil, status.Error(codes.PermissionDenied, "cannot replicate to master")
	}

	res, _ := server.slaveRepl.Replicate(ctx, req)
	return res, nil
}

// func (server *KVStoreServer) HealthCheck(ctx context.Context, req *emptypb.Empty) (*dbv1.HealthCheckResponse, error) {
// 	// Implement basic health check logic
// 	return &dbv1.HealthCheckResponse{Status: dbv1.HealthCheckResponse_SERVING}, nil
// }
