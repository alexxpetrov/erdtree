package components

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/oleksiip-aiola/erdtree/internal/config"
	"github.com/oleksiip-aiola/erdtree/internal/db"
	"github.com/oleksiip-aiola/erdtree/internal/grpc"
	"github.com/oleksiip-aiola/erdtree/internal/server"
	"github.com/oleksiip-aiola/erdtree/internal/wal"
	"github.com/oleksiip-aiola/erdtree/pkg/logger/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Components struct {
	HttpServer *grpc.Server
	Storage    *db.InMemoryDB
	Wal        *wal.WAL
	KVServer   *server.KVStoreServer
}

func InitComponents(cfg *config.Config, logger *slog.Logger, port int, isMaster bool) (*Components, error) {

	var database *db.InMemoryDB
	var kvServer *server.KVStoreServer
	var httpServer *grpc.Server
	var writeAheadLog *wal.WAL
	var err error

	if isMaster {
		walDir := fmt.Sprintf("%s-%d", cfg.Master.WAL.Directory, port)

		writeAheadLog, err = wal.NewWal(walDir, 100*time.Millisecond)
		if err != nil {
			fmt.Errorf("failed  to create wal: %w", err)
		}

		database, err = db.NewInMemoryDb(cfg.Master.Database, writeAheadLog)
		if err != nil {
			fmt.Errorf("failed to create database: %w", err)
		}

		kvServer, err = server.NewMasterKVStoreServer(cfg.Master.WAL.Directory, cfg.Master.SlaveAddresses, database, writeAheadLog)
		if err != nil {
			fmt.Errorf("Failed to create server: %v", err)
		}

		httpServer, err = grpc.New(cfg.Master.Server.Port, kvServer, logger)

		if err != nil {
			return nil, err
		}
	} else {
		walDir := fmt.Sprintf("%s-%d", cfg.Slave.WAL.Directory, port)
		writeAheadLog, err = wal.NewWal(walDir, 100*time.Millisecond)
		if err != nil {
			fmt.Errorf("failed  to create wal: %w", err)
		}

		database, err = db.NewInMemoryDb(cfg.Slave.Database, writeAheadLog)
		if err != nil {

			fmt.Errorf("failed to create database: %w", err)
		}

		kvServer, err = server.NewSlaveKVStoreServer(walDir, database, writeAheadLog)
		if err != nil {

			fmt.Errorf("Failed to create server: %v", err)
		}

		httpServer, err = grpc.New(port, kvServer, logger)

		if err != nil {
			return nil, err
		}
	}

	return &Components{
		HttpServer: httpServer,
		Storage:    database,
		Wal:        writeAheadLog,
		KVServer:   kvServer,
	}, nil
}

func (c *Components) Shutdown() {
	c.HttpServer.Stop()
}

func SetupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slogpretty.SetupPrettySlog()
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
