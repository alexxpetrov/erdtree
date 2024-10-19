# Unified In-Memory Key-Value Database Microservice with Asynchronous Replication

## Project Overview

This project implements an in-memory key-value database microservice with asynchronous replication, designed for use in a chat application. It uses Go, ConnectRPC, Protocol Buffers, and built-in Go libraries for core functionality. The system follows a master-slave architecture with Write-Ahead Logging (WAL) for data persistence and recovery.

## Project Structure

```
kv-database/
├── cmd/
│   └── kvserver/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── db/
│   │   ├── inmemory.go
│   │   └── replication.go
│   ├── server/
│   │   └── server.go
│   ├── api/
│   │   └── kv.proto
│   ├── wal/
│   │   └── wal.go
│   ├── replication/
│   │   ├── master.go
│   │   └── slave.go
│   └── storage/
│       └── storage.go
├── pkg/
│   ├── logger/
│   │   └── logger.go
│   └── metrics/
│       └── metrics.go
├── test/
│   ├── integration/
│   │   └── integration_test.go
│   └── unit/
│       └── unit_test.go
├── deployments/
│   ├── Dockerfile
│   └── docker-compose.yml
├── scripts/
│   ├── build.sh
│   └── test.sh
├── go.mod
├── go.sum
└── README.md
```

## Modules

1. **Database (internal/db)**: Implements the in-memory key-value store and handles asynchronous replication.
2. **Server (internal/server)**: Manages the gRPC server and request handling.
3. **API (internal/api)**: Defines the service API using Protocol Buffers.
4. **Config (internal/config)**: Manages application configuration.
5. **Logger (pkg/logger)**: Provides structured logging for the application.
6. **Metrics (pkg/metrics)**: Implements metrics collection and reporting.
7. **WAL (internal/wal)**: Implements Write-Ahead Logging for data persistence and recovery.
8. **Replication (internal/replication)**: Handles master-slave replication logic.
9. **Storage (internal/storage)**: Manages the storage layer, including in-memory and disk operations.

## Architecture

The key-value database follows a master-slave architecture with the following key components:

- **Compute Layer**: Handles incoming requests and basic operations.
- **Storage Layer**: Manages data storage, including in-memory operations and disk persistence.
- **In-Memory Engine**: Provides fast access to data stored in memory.
- **WAL (Write-Ahead Log)**: Ensures data durability and supports recovery.
- **Async Replication**: Implements master-slave replication for data redundancy.

Both master and slave nodes have similar structures, with the master handling write operations and coordinating replication.

## Step-by-Step Implementation Guide

1. **Project Setup**
   - Initialize the Go module:
     ```
     go mod init github.com/oleksiip-aiola/erdtree
     ```
   - Set up the project structure as outlined above.
   - Create a basic `README.md` with project description and structure.

   Link: [Go Modules](https://go.dev/blog/using-go-modules)

2. **Define the API (internal/api/kv.proto)**
   - Create a Protocol Buffer definition for your key-value store operations.
   - Include methods for Get, Set, Delete, and ReplicationSync.

   ```protobuf
   syntax = "proto3";

   package kv.v1;

   option go_package = "github.com/oleksiip-aiola/erdtree/internal/api;api";

   service KVStore {
     rpc Get(GetRequest) returns (GetResponse) {}
     rpc Set(SetRequest) returns (SetResponse) {}
     rpc Delete(DeleteRequest) returns (DeleteResponse) {}
     rpc ReplicationSync(ReplicationSyncRequest) returns (ReplicationSyncResponse) {}
   }

   // Define message types for requests and responses
   ```

   Link: [Protocol Buffers](https://protobuf.dev/overview/)

3. **Implement In-Memory Database (internal/db/inmemory.go)**
   - Use `sync.Map` for thread-safe operations.
   - Implement Get, Set, and Delete methods.

   ```go
   type InMemoryDB struct {
       data sync.Map
   }

   func (db *InMemoryDB) Get(key string) (interface{}, bool) {
       return db.data.Load(key)
   }

   func (db *InMemoryDB) Set(key string, value interface{}) {
       db.data.Store(key, value)
   }

   func (db *InMemoryDB) Delete(key string) {
       db.data.Delete(key)
   }
   ```

   Link: [sync.Map](https://pkg.go.dev/sync#Map)

4. **Implement Write-Ahead Logging (internal/wal/wal.go)**
   - Create a WAL system for data persistence and recovery.

   ```go
   type WAL struct {
       file *os.File
       mu   sync.Mutex
   }

   func (w *WAL) AppendEntry(entry []byte) error {
       w.mu.Lock()
       defer w.mu.Unlock()
       _, err := w.file.Write(append(entry, '\n'))
       return err
   }

   func (w *WAL) Recover() ([][]byte, error) {
       // Implement recovery logic
   }
   ```

   Link: [Go File I/O](https://pkg.go.dev/os#File)

5. **Implement Storage Layer (internal/storage/storage.go)**
   - Integrate WAL with in-memory storage.

   ```go
   type Storage struct {
       data map[string][]byte
       wal  *wal.WAL
       mu   sync.RWMutex
   }

   func (s *Storage) Set(key string, value []byte) error {
       s.mu.Lock()
       defer s.mu.Unlock()
       s.data[key] = value
       return s.wal.AppendEntry([]byte(fmt.Sprintf("SET %s %s", key, value)))
   }

   // Implement Get, Delete, and recovery methods
   ```

6. **Implement Master Replication (internal/replication/master.go)**
   - Handle replication logic on the master node.

   ```go
   type MasterReplication struct {
       slaves []string
       wal    *wal.WAL
   }

   func (mr *MasterReplication) ReplicateToSlaves(entry []byte) {
       for _, slave := range mr.slaves {
           go mr.sendToSlave(slave, entry)
       }
   }

   func (mr *MasterReplication) sendToSlave(slaveAddr string, entry []byte) {
       // Implement logic to send entry to slave
   }
   ```

7. **Implement Slave Replication (internal/replication/slave.go)**
   - Handle replication logic on slave nodes.

   ```go
   type SlaveReplication struct {
       masterAddr string
       storage    *storage.Storage
   }

   func (sr *SlaveReplication) PullFromMaster() {
       // Implement logic to pull updates from master
   }

   func (sr *SlaveReplication) ApplyUpdate(entry []byte) error {
       // Apply the update to local storage
       return sr.storage.ApplyEntry(entry)
   }
   ```

8. **Implement gRPC Server (internal/server/server.go)**
   - Use ConnectRPC to implement the server.
   - Handle Get, Set, Delete, and ReplicationSync operations.

   ```go
   type KVServer struct {
       storage   *storage.Storage
       replication interface {
           ReplicateToSlaves(entry []byte)
       }
   }

   func (s *KVServer) Set(ctx context.Context, req *api.SetRequest) (*api.SetResponse, error) {
       err := s.storage.Set(req.Key, req.Value)
       if err != nil {
           return nil, status.Error(codes.Internal, "Failed to set value")
       }
       s.replication.ReplicateToSlaves([]byte(fmt.Sprintf("SET %s %s", req.Key, req.Value)))
       return &api.SetResponse{}, nil
   }

   // Implement Get, Delete, and ReplicationSync similarly
   ```

   Link: [ConnectRPC](https://connectrpc.com/docs/go/getting-started/)

9. **Configuration Management (internal/config/config.go)**
   - Use `github.com/spf13/viper` for configuration management.
   - Support both file-based and environment variable configuration.

   ```go
   type Config struct {
       ServerPort     int
       PeerAddresses  []string
       Role           string // "master" or "slave"
       MasterAddress  string // Used by slaves to connect to master
   }

   func LoadConfig() (*Config, error) {
       viper.SetConfigName("config")
       viper.AddConfigPath(".")
       viper.AutomaticEnv()

       if err := viper.ReadInConfig(); err != nil {
           return nil, err
       }

       var config Config
       err := viper.Unmarshal(&config)
       return &config, err
   }
   ```

   Link: [Viper](https://github.com/spf13/viper)

10. **Logging (pkg/logger/logger.go)**
    - Use `go.uber.org/zap`/`slog` for structured logging.

    ```go
    var logger *zap.Logger

    func InitLogger() {
        var err error
        logger, err = zap.NewProduction()
        if err != nil {
            log.Fatalf("Can't initialize zap logger: %v", err)
        }
    }

    func Info(msg string, fields ...zap.Field) {
        logger.Info(msg, fields...)
    }

    // Implement Debug, Error, etc.
    ```

    Link: [Zap](https://pkg.go.dev/go.uber.org/zap)

11. **Metrics (pkg/metrics/metrics.go)**
    - Use `github.com/prometheus/client_golang` for metrics collection.

    ```go
    var (
        OperationsTotal = promauto.NewCounterVec(
            prometheus.CounterOpts{
                Name: "kv_operations_total",
                Help: "The total number of processed operations",
            },
            []string{"operation"},
        )
    )

    func IncrementOperations(operation string) {
        OperationsTotal.WithLabelValues(operation).Inc()
    }
    ```

    Link: [Prometheus Client Go](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus)

12. **Main Application (cmd/kvserver/main.go)**
    - Initialize all components and start the server.

    ```go
    func main() {
        config, err := config.LoadConfig()
        if err != nil {
            log.Fatalf("Failed to load config: %v", err)
        }

        logger.InitLogger()
        defer logger.Sync()

        storage := storage.NewStorage()
        var replication interface{}
        if config.Role == "master" {
            replication = replication.NewMasterReplication(config.SlaveAddresses, storage.WAL())
        } else {
            replication = replication.NewSlaveReplication(config.MasterAddress, storage)
        }

        server := server.NewKVServer(storage, replication)

        // Start the gRPC server
        lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ServerPort))
        if err != nil {
            logger.Fatal("Failed to listen", zap.Error(err))
        }

        grpcServer := grpc.NewServer()
        api.RegisterKVStoreServer(grpcServer, server)

        logger.Info("Starting gRPC server", zap.Int("port", config.ServerPort))
        if err := grpcServer.Serve(lis); err != nil {
            logger.Fatal("Failed to serve", zap.Error(err))
        }
    }
    ```

13. **Testing (test/)**
    - Implement unit tests for individual components.
    - Create integration tests for the entire service.

    ```go
    func TestInMemoryDB(t *testing.T) {
        db := NewInMemoryDB()
        db.Set("key", "value")
        value, ok := db.Get("key")
        assert.True(t, ok)
        assert.Equal(t, "value", value)
    }

    // Implement more tests for other components and integration scenarios
    ```

    Link: [Go Testing](https://pkg.go.dev/testing)

14. **Dockerization (deployments/Dockerfile)**
    - Create a multi-stage Dockerfile for efficient builds.

    ```dockerfile
    # Build stage
    FROM golang:1.17 AS build
    WORKDIR /app
    COPY . .
    RUN CGO_ENABLED=0 GOOS=linux go build -o kvserver ./cmd/kvserver

    # Final stage
    FROM alpine:latest
    RUN apk --no-cache add ca-certificates
    WORKDIR /root/
    COPY --from=build /app/kvserver .
    COPY --from=build /app/config.yaml .
    CMD ["./kvserver"]
    ```

    Link: [Docker Multi-stage Builds](https://docs.docker.com/build/building/multi-stage/)

15. **CI/CD Setup**
    - Create a GitHub Actions workflow for continuous integration.
    - Implement automatic testing and deployment.

    ```yaml
    name: CI/CD

    on:
      push:
        branches: [ main ]
      pull_request:
        branches: [ main ]

    jobs:
      build:
        runs-on: ubuntu-latest
        steps:
        - uses: actions/checkout@v2
        - name: Set up Go
          uses: actions/setup-go@v2
          with:
            go-version: 1.17
        - name: Build
          run: go build -v ./...
        - name: Test
          run: go test -v ./...

      deploy:
        needs: build
        runs-on: ubuntu-latest
        steps:
        - name: Deploy to production
          # Add your deployment steps here
    ```

    Link: [GitHub Actions](https://docs.github.com/en/actions)

16. **Documentation**
    - Update the README.md with detailed information about the project, its structure, and how to run it.
    - Add inline comments and generate API documentation.

    ```go
    // GenerateAPIDoc generates API documentation using protoc-gen-doc
    //go:generate protoc --doc_out=. --doc_opt=markdown,api.md api/kv.proto
    ```

    Link: [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc)

17. **Performance Tuning**
    - Use Go's built-in pprof for profiling.
    - Implement benchmarks for critical paths.

    ```go
    import _ "net/http/pprof"

    func main() {
        go func() {
            log.Println(http.ListenAndServe("localhost:6060", nil))
        }()
        // Rest of your main function
    }
    ```

    Link: [pprof](https://pkg.go.dev/net/http/pprof)

18. **Error Handling and Logging**
    - Implement centralized error handling.
    - Ensure all errors are properly logged and monitored.

    ```go
    func handleError(err error, msg string) {
        logger.Error(msg, zap.Error(err))
        // Implement error reporting to a centralized system if needed
    }
    ```

19. **Implement Health Checks**
    - Add health check endpoints for both the master and slave nodes.
    - Include checks for replication lag and system resources.

    ```go
    func (s *KVServer) HealthCheck(ctx context.Context, req *api.HealthCheckRequest) (*api.HealthCheckResponse, error) {
        status := api.HealthCheckResponse_SERVING
        if !s.isHealthy() {
            status = api.HealthCheckResponse_NOT_SERVING
        }
        return &api.HealthCheckResponse{Status: status}, nil
    }

    func (s *KVServer) isHealthy() bool {
        // Implement health check logic
        return true
    }
    ```

    Link: [gRPC Health Checking Protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)

20. **Implement Graceful Shutdown**
    - Ensure the server can shut down gracefully, completing in-flight requests.

    ```go
    func main() {
        // ... (previous main function code)

        go func() {
            sigChan := make(chan os.Signal, 1)
            signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
            <-sigChan
            logger.Info("Received shutdown signal, initiating graceful shutdown")
            grpcServer.GracefulStop()
        }()

        if err := grpcServer.Serve(lis); err != nil {
            logger.Fatal("Failed to serve", zap.Error(err))
        }
    }
    ```

21. **Implement Rate Limiting**
    - Add rate limiting to protect the service from overwhelming traffic.

    ```go
    import "golang.org/x/time/rate"

    type KVServer struct {
        // ... (previous fields)
        limiter *rate.Limiter
    }

    func (s *KVServer) Set(ctx context.Context, req *api.SetRequest) (*api.SetResponse, error) {
        if !s.limiter.Allow() {
            return nil, status.Error(codes.ResourceExhausted, "Rate limit exceeded")
        }
        // ... (rest of the Set method)
    }
    ```

    Link: [rate package](https://pkg.go.dev/golang.org/x/time/rate)

22. **Implement Monitoring and Alerting**
    - Set up monitoring for key metrics (e.g., request rate, error rate, replication lag).
    - Configure alerts for critical issues.

    ```go
    // In your metrics package
    var (
        ReplicationLag = promauto.NewGauge(prometheus.GaugeOpts{
            Name: "kv_replication_lag_seconds",
            Help: "The current replication lag in seconds",
        })
    )

    // In your replication logic
    func updateReplicationLag(lag time.Duration) {
        ReplicationLag.Set(lag.Seconds())
    }
    ```

23. **Implement Data Compression**
    - Add optional compression for stored values to reduce memory usage.

    ```go
    import "compress/gzip"

    func compressValue(value []byte) ([]byte, error) {
        var b bytes.Buffer
        w := gzip.NewWriter(&b)
        _, err := w.Write(value)
        w.Close()
        if err != nil {
            return nil, err
        }
        return b.Bytes(), nil
    }

    func decompressValue(compressedValue []byte) ([]byte, error) {
        b := bytes.NewReader(compressedValue)
        r, err := gzip.NewReader(b)
        if err != nil {
            return nil, err
        }
        defer r.Close()
        return ioutil.ReadAll(r)
    }
    ```

24. **Implement Periodic Consistency Checks**
    - Add a background job to periodically verify data consistency between master and slaves.

    ```go
    func (mr *MasterReplication) StartConsistencyChecks(interval time.Duration) {
        ticker := time.NewTicker(interval)
        go func() {
            for range ticker.C {
                mr.performConsistencyCheck()
            }
        }()
    }

    func (mr *MasterReplication) performConsistencyCheck() {
        // Implement consistency check logic
    }
    ```

25. **Optimize for Production**
    - Fine-tune Go garbage collection settings for production workloads.
    - Implement connection pooling for slave connections.

    ```go
    import "runtime/debug"

    func init() {
        debug.SetGCPercent(100)  // Adjust based on your workload
    }
    ```

26. **Security Considerations**
    - Implement TLS for all network communications.
    - Add authentication and authorization mechanisms.

    ```go
    // In your main function
    creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
    if err != nil {
        log.Fatalf("Failed to load TLS keys: %v", err)
    }
    grpcServer := grpc.NewServer(grpc.Creds(creds))
    ```

    Link: [gRPC Authentication](https://grpc.io/docs/guides/auth/)

27. **Disaster Recovery Planning**
    - Implement regular backups of the WAL and in-memory state.
    - Create a disaster recovery plan and test it regularly.

    ```go
    func (s *Storage) Backup() error {
        // Implement backup logic
    }

    func (s *Storage) Restore(backupFile string) error {
        // Implement restore logic
    }
    ```

28. **Documentation and Runbooks**
    - Create comprehensive documentation for operating the service.
    - Develop runbooks for common operational tasks and incident response.

    ```markdown
    # Operational Runbook

    ## Handling High Replication Lag
    1. Check network connectivity between master and slaves.
    2. Verify master is not overloaded.
    3. If necessary, restart the replication process on affected slaves.

    ## Performing a Failover
    1. Identify the slave with the most up-to-date data.
    2. Stop writes to the current master.
    3. Promote the identified slave to master.
    4. Update configuration on all nodes.
    5. Restart the service with new roles.

    ...
    ```

This completes the unified project plan for the in-memory key-value database microservice with asynchronous replication. The plan covers all aspects of development, from initial setup to production optimization and operational considerations. By following this plan, you'll create a robust, scalable, and maintainable solution that meets the requirements of a production-ready microservice for your chat application.

Remember to iterate on this plan as you develop, adjusting based on specific requirements, performance metrics, and operational feedback. Regular code reviews, performance testing, and security audits should be part of your ongoing development process.
