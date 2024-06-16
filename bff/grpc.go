package main

import (
	"fmt"
	"net"

	"github.com/gulldan/lct2024_copyright/bff/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	grpchandler "github.com/gulldan/lct2024_copyright/bff/internal/transport/grpc"
	bffv1 "github.com/gulldan/lct2024_copyright/bff/proto/bff/v1"
)

func startGRPCServer(cfg *config.Config, logger *zerolog.Logger, grpcEndpoint string) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		return nil, fmt.Errorf("can't start listener for grpc server: %w", err)
	}

	var opts []grpc.ServerOption

	middlewares := []grpc.UnaryServerInterceptor{}

	opts = append(opts, grpc.ChainUnaryInterceptor(middlewares...))

	grpcServer := grpc.NewServer(opts...)

	hdl := grpchandler.New(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("grpc handler init failed: %w", err)
	}

	bffv1.RegisterScanTasksServiceServer(grpcServer, hdl)

	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			logger.Error().Err(err).Msg("can't start grpc server")
		}
	}()

	log.Info().Str("addr", grpcEndpoint).Msg("started grpc server")

	return grpcServer, nil
}
