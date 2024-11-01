# Distributed Key-Value Store System

## System Overview
This is a distributed key-value store system implementing a master-slave replication pattern with Write-Ahead Logging (WAL) for durability. The system is built using Go, Connect-RPC, and Protocol Buffers, designed to provide high availability and eventual consistency.

## Architecture
![alt text](image.png)

## Flow Diagram 
Master Node                                 Slave Node
+------------------------+                +-----------------------+
|    write requests      |                |     read requests     |
|          ↓             |                |          ↓            |
|   [Compute Layer]      |                |   [Compute Layer]     |
|          ↓             |      push      |          ↓            |
|   [In-Memory Engine]   |   --------->   |   [In-Memory Engine]  |
|          ↓             |                |          ↓            |
|   [Storage Layer]      |                |     [Storage]         |
|          ↓             |                |          ↓            |
|        [WAL]           |                |        [WAL]          |
|          ↓             |                |          ↓            |
|       [DISK]           |                |       [DISK]          |
+------------------------+                +-----------------------+

The system consists of several key components:
- Master Node
- Slave Nodes
- Write-Ahead Log (WAL)
- In-Memory Storage Engine
- Replication System

### System Design
```
Master Node
├── In-Memory Storage (sync.Map)
├── Write-Ahead Log (WAL)
└── Compute Layer
    └── TCP Server (Connect-RPC)

Slave Node
├── In-Memory Storage (sync.Map)
├── Write-Ahead Log (WAL)
└── Compute Layer
    └── TCP Server (Connect-RPC)
```

## Core Components

### 1. Master Node
**Purpose**: Handles all write operations and coordinates replication to slave nodes.

**Key Responsibilities**:
- Processes write operations (SET, DELETE)
- Maintains the primary copy of data
- Logs operations to WAL
- Keeps an ongoing syncloop that syncs with slaves on an interval

**Location**: `internal/master/master.go`

**Key Methods**:
- `AddSlave(addr)`: Add a slave connection
- `RemoveSlave(addr)`: Remove a slave connection
- `ReplicateEntry(entry)`: Manual entry replication

### 2. Slave Node
**Purpose**: Maintains a replicated copy of the data and handles read operations.

**Key Responsibilities**:
- Handles read operations
- Periodically pulls updates from master
- Maintains local WAL
- Provides read-only access to data

**Location**: `internal/slave/slave.go`

**Key Methods**:
- `Get(key)`: Read operations
- `Replicate(entries)`: Receives entries that are up to be replicated into current slave node

### 3. KV Store Server (API)
**Purpose**: Provides methods to be used by nodes to manage in memory data.

**Location**: `internal/server/server.go`

**Key Features**:
- Atomic data management using sync.Map
# Why sync.Map over a regular Map? Because we're creating a cache storage that cannot be edited and can only be READ/DELETED/VIEWED. 
# It's one of 2 cases when a sync.Map is more optimized, than a regular map. 
# https://pkg.go.dev/sync#Map

**Key Methods**:
- `Set(key, value)`: Write operations
- `Delete(key)`: Delete operations
- `Get(key)`: Get record by key

### 4. Write-Ahead Log (WAL)
**Purpose**: Ensures durability of operations and enables recovery after crashes.

**Location**: `internal/wal/wal.go`

**Key Features**:
- Sequential write operations
- Entry format includes timestamps and operation types
- Supports recovery operations
- Manages log file rotation

**Key Methods**:
- `AppendEntry(entry)`: Logs new operations
- `Recover()`: Rebuilds state from log
- `GetEntriesSince(timestamp)`: Retrieves entries for replication

### 5. Storage Engine
**Purpose**: Manages in-memory data storage using sync.Map for thread-safe operations.

**Location**: `internal/storage/storage.go`

**Key Features**:
- Thread-safe operations using sync.Map
- Support for TTL (Time-To-Live)
- Memory management and eviction

## Communication Protocol

### Connect-RPC Service Definition
**Location**: `api/kvstore.proto`

```protobuf
service KVStore {
    rpc Get(GetRequest) returns (GetResponse);
    rpc Set(SetRequest) returns (SetResponse);
    rpc Delete(DeleteRequest) returns (DeleteResponse);
    rpc GetUpdates(GetUpdatesRequest) returns (GetUpdatesResponse);
}
```

### Data Flow
1. **Write Operation**:
   ```
   Client -> Master Node -> WAL -> In-Memory Storage
                        -> Slave Nodes (via push-based replication)
   ```

2. **Read Operation**:
   ```
   Client -> Any Node (Slave) -> In-Memory Storage
   ```

3. **Replication Flow**:
   ```
   Master -> Slave (Push updates)
   Slave -> Master (Updates Response)
   Slave -> Local WAL -> In Memory DB
   ```

## Configuration

### Master Node Configuration
```yaml
server:
  port: 8080
  address: "0.0.0.0"

wal:
  directory: "/var/log/kvstore/master"
  sync_interval: "1s"
  max_file_size: "100MB"

storage:
  max_entries: 1000000
  eviction_policy: "lru"
```

### Slave Node Configuration
```yaml
server:
  port: 8081
  address: "0.0.0.0"

master:
  address: "master:8080"
  sync_interval: "1s"

wal:
  directory: "/var/log/kvstore/slave"
  sync_interval: "1s"
```

## Getting Started

### Prerequisites
- Go 1.19 or later
- Protocol Buffer compiler
- Connect-go

### Building the System
```bash
# Generate Protocol Buffer code
protoc -I . \
    --go_out=. --go_opt=paths=source_relative \
    --connect-go_out=. --connect-go_opt=paths=source_relative \
    api/kvstore.proto

# Build master and slave binaries
go build -o master ./cmd/master
go build -o slave ./cmd/slave
```

### Running the System
1. Start the master node:
```bash
./master -config master-config.yaml
```

2. Start slave nodes:
```bash
./slave -config slave-config.yaml
```

## Recovery Process

### Master Recovery
1. Reads WAL entries on startup
2. Rebuilds in-memory state
3. Resumes normal operation

### Slave Recovery
1. Reads local WAL entries
2. Rebuilds local state
3. Syncs with master for missing updates

## Monitoring and Maintenance

### Key Metrics
- Replication lag
- Operation latency
- Memory usage
- WAL size
- Error rates

### Health Checks
- Master node availability
- Slave node status
- Replication status
- WAL write success rate

## Best Practices

### Production Deployment
1. Use secure connections (TLS)
2. Implement proper monitoring
3. Regular backup of WAL files
4. Monitor disk space for WAL
5. Implement proper error handling
6. Set up alerting for critical issues

### Performance Optimization
1. Tune WAL sync interval
2. Adjust batch sizes for replication
3. Configure appropriate TTL values
4. Monitor and adjust memory usage
5. Implement proper connection pooling

## Troubleshooting

### Common Issues
1. Replication lag
   - Check network connectivity
   - Verify slave sync interval
   - Monitor master load

2. High latency
   - Check WAL sync frequency
   - Monitor memory usage
   - Verify network conditions

3. Data inconsistency
   - Check WAL integrity
   - Verify replication status
   - Review error logs

## Security Considerations
1. Network security
2. Authentication
3. Authorization
4. Data encryption
5. Secure configuration management

## Contributing
Please read CONTRIBUTING.md for details on our code of conduct and the process for submitting pull requests.

## License
This project is licensed under the MIT License - see the LICENSE file for details.