package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	"repo-stat/processor/config"
	"repo-stat/processor/internal/adapter/collector"
	grpccontroller "repo-stat/processor/internal/controller/grpc"
	"repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/processor"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting processor server...")
	log.Debug("debug messages are enabled")

	collectorClient, err := collector.NewClient(cfg.Services.Collector)
	if err != nil {
		return fmt.Errorf("cannot connect to collector: %w", err)
	}
	log.Info("successfully connected to collector", "address", cfg.Services.Collector)

	repoUseCase := usecase.NewService(collectorClient)
	processorServer := grpccontroller.NewServer(repoUseCase)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	processorpb.RegisterProcessorServiceServer(srv.GRPC(), processorServer)

	log.Info("processor is listening...", "address", cfg.GRPC.Address)
	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	if err := run(ctx); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Printf("launching server error: %s\n", err)
		}
		cancel()
		os.Exit(1)
	}
	cancel()
}
