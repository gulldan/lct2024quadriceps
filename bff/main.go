package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gulldan/lct2024_copyright/bff/pkg/config"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

//go:embed swagger-ui/docs
var swaggerDocsFS embed.FS

func main() {
	cfg, logLevel, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	log := zerolog.New(os.Stdout).Level(*logLevel).With().Timestamp().Logger()

	grpcEndpoint := "localhost" + cfg.Grpc.Address

	grpcServer, err := startGRPCServer(cfg, &log, grpcEndpoint)
	if err != nil {
		panic(err)
	}

	httpServer, err := startHTTPServer(cfg, &log, grpcEndpoint)
	if err != nil {
		log.Error().Err(err).Msg("start http server failed")
	}

	if err := gracefulShutdown(&log, httpServer, grpcServer); err != nil {
		log.Error().Err(err).Msg("graceful shutdown failed")
	}
}

func gracefulShutdown(logger *zerolog.Logger, h *http.Server, g *grpc.Server) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)

	sig := <-sigs
	logger.Info().Str("signal", sig.String()).Msg("signal received, graceful shutdown")

	defer cancel()

	if err := h.Shutdown(ctx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	g.Stop()

	return nil
}
