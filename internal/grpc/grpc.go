package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/oleksiip-aiola/erdtree/gen/api/v1/dbv1connect"
	"github.com/oleksiip-aiola/erdtree/internal/server"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	mux             *http.ServeMux
	path            string
	corsHandler     http.Handler
	srv             *http.Server
	shutDownTimeout time.Duration
	logger          *slog.Logger
}

func New(port int, kvStore *server.KVStoreServer, logger *slog.Logger) (*Server, error) {
	mux := http.NewServeMux()

	path, handler := dbv1connect.NewErdtreeStoreHandler(kvStore)

	fmt.Printf("ConnectRPC is serving at :%d\n", port)

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

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}
	shutDownTimeout := 10 * time.Second

	return &Server{
		mux,
		path,
		corsHandler,
		srv,
		shutDownTimeout,
		logger,
	}, nil
}

func (server *Server) Run(ctx context.Context) error {
	errResult := make(chan error)
	go func() {
		server.mux.Handle(server.path, server.corsHandler)

		server.logger.Info(fmt.Sprintf("starting listening: %s", server.srv.Addr))

		// if server.certFilePath != "" && server.keyFilePath != "" {
		// 	errResult <- server.srv.ListenAndServeTLS(server.certFilePath, server.keyFilePath)
		// }
		server.srv.ListenAndServe()
	}()

	var err error
	select {
	case <-ctx.Done():
		return ctx.Err()

	case err = <-errResult:
	}
	return err
}

func (server *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), server.shutDownTimeout)
	defer cancel()
	err := server.srv.Shutdown(ctx)
	if err != nil {
		server.logger.Error("failed to shutdown HTTP Server", slog.String("error", err.Error()))
	}
}
