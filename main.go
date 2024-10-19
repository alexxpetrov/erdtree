package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/oleksiip-aiola/erdtree/gen/api/v1/dbv1connect"
	"github.com/oleksiip-aiola/erdtree/internal/config"
	"github.com/oleksiip-aiola/erdtree/internal/db"
	"github.com/oleksiip-aiola/erdtree/internal/server"
	"github.com/oleksiip-aiola/erdtree/internal/wal"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Errorf("Failed to load config: %v", err)
	}

	if err := config.ValidateConfig(cfg); err != nil {
		fmt.Errorf("Invalid configuration: %v", err)
	}

	wal, err := wal.NewWal("walstore", 100*time.Millisecond)
	if err != nil {
		fmt.Errorf("failed  to create wal: %w", err)
	}
	database, err := db.NewInMemoryDb(cfg.Database, wal)
	if err != nil {
		fmt.Errorf("failed to create database: %w", err)
	}

	master, err := server.NewKVStoreServer(cfg.WAL.Directory, cfg.IsMaster(), cfg.Master.SlaveAddresses, database, wal)
	if err != nil {
		fmt.Errorf("Failed to create server: %v", err)
	}

	mux := http.NewServeMux()

	// interceptors := connect.WithInterceptors(NewAuthInterceptor())
	path, handler := dbv1connect.NewErdtreeStoreHandler(master)

	fmt.Printf("ConnectRPC is serving at :%s\n", os.Getenv("PORT"))

	publicUrl := os.Getenv("PUBLIC_URL")

	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if publicUrl != "" {
			w.Header().Set("Access-Control-Allow-Origin", publicUrl)

		} else {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Set-Cookie, connect-protocol-version")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			return
		}

		// Call the next handler with the updated context
		handler.ServeHTTP(w, r)
	})

	mux.Handle(path, corsHandler)

	http.ListenAndServe(
		":"+os.Getenv("PORT"),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
