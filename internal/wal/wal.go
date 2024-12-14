package wal

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	dbv1 "github.com/oleksiip-aiola/erdtree/gen/api/v1"
	"google.golang.org/protobuf/proto"
)

const (
	logFilePrefix  = "wal-"
	logFilePostfix = ".log"
	maxLogSize     = 1024 * 1024 * 100 // 100MB
)

type LogEntry struct {
	Timestamp int64
	Operation dbv1.Operation
	Key       string
	Value     []byte
	ExpiresAt int64 // Unix nano
}

type WAL struct {
	dir            string
	currentFile    *os.File
	currentSize    int64
	mx             sync.Mutex
	bufferedWriter *bufio.Writer
	syncInterval   time.Duration
	lastSync       time.Time
}

func NewWal(dir string, syncInterval time.Duration) (*WAL, error) {
	absPath, _ := filepath.Abs(dir)
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create WAL directory: %w", err)
	}

	wal := &WAL{
		dir:          dir,
		syncInterval: syncInterval,
	}

	if err := wal.openNewLogFile(); err != nil {
		return nil, fmt.Errorf("failed to open log file %w", err)
	}

	return wal, nil
}

func (wal *WAL) openNewLogFile() error {
	wal.mx.Lock()
	defer wal.mx.Unlock()

	if wal.currentFile != nil {
		wal.bufferedWriter.Flush()
		wal.currentFile.Close()
	}

	nowUtc := time.Now().UTC()
	fileName := fmt.Sprintf("%s-%d-%d-%d-%s", logFilePrefix, nowUtc.Month(), nowUtc.Day(), nowUtc.Hour(), logFilePostfix)
	filePath := filepath.Join(wal.dir, fileName)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create/open file: %w", err)
	}

	wal.currentFile = file
	wal.currentSize = 0
	wal.bufferedWriter = bufio.NewWriter(file)
	wal.lastSync = nowUtc

	return nil
}

func (wal *WAL) AppendEntry(entry *dbv1.LogEntry) error {
	wal.mx.Lock()
	defer wal.mx.Unlock()

	data, err := proto.Marshal(entry)

	if err != nil {
		return fmt.Errorf("failed to marshal log entry %w", err)
	}

	entryLength := uint32(len(data))

	if err := binary.Write(wal.bufferedWriter, binary.BigEndian, entryLength); err != nil {
		return fmt.Errorf("failed to write entry length %w", err)
	}

	if _, err := wal.bufferedWriter.Write(data); err != nil {
		return fmt.Errorf("failed to write buffer data %w", err)
	}
	wal.currentSize += int64(4 + len(data))

	if time.Since(wal.lastSync) >= wal.syncInterval {
		if err := wal.sync(); err != nil {
			return fmt.Errorf("failed to sync wal %w", err)
		}
	}

	if wal.currentSize >= maxLogSize {
		if err := wal.openNewLogFile(); err != nil {
			return fmt.Errorf("failed to open new log file: %w", err)
		}
	}
	return nil
}

func (wal *WAL) sync() error {
	if err := wal.bufferedWriter.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffered writer: %w", err)
	}

	if err := wal.currentFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync current file: %w", err)
	}

	wal.lastSync = time.Now()

	return nil
}

func (wal *WAL) Flush() error {
	wal.mx.Lock()
	defer wal.mx.Unlock()

	return wal.sync()
}

func (wal *WAL) Close() error {
	wal.mx.Lock()
	defer wal.mx.Unlock()

	if err := wal.sync(); err != nil {
		return fmt.Errorf("failed to sync file before closing: %w", err)
	}

	return wal.currentFile.Close()
}

func (w *WAL) Recover(timestamp int64) ([]*dbv1.LogEntry, error) {
	w.mx.Lock()
	defer w.mx.Unlock()

	var entries []*dbv1.LogEntry

	nowUtc := time.Now().UTC()

	if timestamp != 0 {
		nowUtc = time.Unix(timestamp, 0).UTC()
	}

	files, err := filepath.Glob(filepath.Join(w.dir, fmt.Sprintf("%s-%d-%d-%d-%s", logFilePrefix, nowUtc.Month(), nowUtc.Day(), nowUtc.Hour(), logFilePostfix)))

	if err != nil {
		return nil, fmt.Errorf("failed to gather files: %w", err)
	}
	for _, file := range files {
		fileEntries, err := w.recoverFromFile(file)

		if err != nil {
			return nil, err
		}

		entries = append(entries, fileEntries...)
	}

	return entries, nil
}

func (wal *WAL) recoverFromFile(filePath string) ([]*dbv1.LogEntry, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var entries []*dbv1.LogEntry
	reader := bufio.NewReader(file)

	for {
		var entryLength int32
		if err := binary.Read(reader, binary.BigEndian, &entryLength); err != nil {
			if err.Error() == io.EOF.Error() {
				break
			}
			return nil, fmt.Errorf("failed to read entry length: %w", err)
		}

		data := make([]byte, entryLength)

		if _, err := reader.Read(data); err != nil {
			return nil, fmt.Errorf("failed to read entry data: %w", err)
		}

		var entry dbv1.LogEntry
		if err := proto.Unmarshal(data, &entry); err != nil {
			return nil, fmt.Errorf("failed to unmarshal data: %w", err)
		}

		entries = append(entries, &entry)
	}

	return entries, nil
}

func (wal *WAL) GetEntriesSince(lastSync int64) ([]*dbv1.LogEntry, error) {
	entries, err := wal.Recover(lastSync)

	if err != nil {
		return nil, err
	}

	return entries, err
}
