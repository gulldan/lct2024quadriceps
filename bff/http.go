package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gulldan/lct2024_copyright/bff/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	gateway "github.com/gulldan/lct2024_copyright/bff/internal/transport/grpc_gateway"
)

func startHTTPServer(cfg *config.Config, logger *zerolog.Logger, grpcEndpoint string) (*http.Server, error) {
	mux, err := gateway.New(logger, grpcEndpoint)
	if err != nil {
		return nil, fmt.Errorf("can't get http gateway: %w", err)
	}

	mux, err = gateway.AddSwagger(swaggerDocsFS, mux)
	if err != nil {
		return nil, fmt.Errorf("failed to add swagger-ui: %w", err)
	}

	mux = gateway.AddCors(mux)
	mux = gateway.HandleBinaryFileUpload(mux, cfg, logger)
	mux = gateway.HandleOriginalFileUpload(mux, cfg, logger)

	httpServer := http.Server{
		Addr:              "0.0.0.0:" + cfg.HTTPPort,
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		IdleTimeout:       10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error().Err(err).Msg("can't start http server")
		}
	}()

	log.Info().Str("addr", "0.0.0.0:"+cfg.HTTPPort).Msg("started http server")

	return &httpServer, nil
}
