package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/oleksiip-aiola/erdtree/internal/components"
	"github.com/oleksiip-aiola/erdtree/internal/config"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg, err := config.LoadConfig("../config.yaml")
	if err != nil {
		fmt.Errorf("Failed to load config: %v", err)
	}

	if err := config.ValidateConfig(cfg); err != nil {
		fmt.Errorf("Invalid configuration: %v", err)
	}

	logger := components.SetupLogger(cfg.Env)
	for _, address := range cfg.Master.SlaveAddresses {
		address := strings.Split(address, ":")
		port, _ := strconv.Atoi(address[1])
		go createNode(cfg, logger, port, false)
	}

	createNode(cfg, logger, cfg.Master.Server.Port, true)
}

func createNode(cfg *config.Config, logger *slog.Logger, port int, isMaster bool) {
	var serverComponents *components.Components
	var err error

	serverComponents, err = components.InitComponents(cfg, logger, port, isMaster)

	if err != nil {
		logger.Error("bad configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer serverComponents.Shutdown()

	eg, ctx := errgroup.WithContext(context.Background())
	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		return serverComponents.HttpServer.Run(ctx)
	})

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case s := <-sigQuit:
			logger.Info("Captured signal", slog.String("signal", s.String()))
			return fmt.Errorf("captured signal: %v", s)
		}
	})

	err = eg.Wait()
	logger.Info("Gracefully shutting down the servers", slog.String("error", err.Error()))

}
