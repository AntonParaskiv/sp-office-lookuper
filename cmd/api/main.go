package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"sp-office-lookuper/internal/config"
	"sp-office-lookuper/internal/grpcserver"
	"sp-office-lookuper/internal/httpserver"
	"sp-office-lookuper/internal/logging"
	"sp-office-lookuper/internal/storage"
	"sp-office-lookuper/internal/tracer"

	"github.com/joho/godotenv"
)

//nolint:funlen // it's ok
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}
	os.Setenv("SERVICE_NAME", "sp-office-lookuper")

	cfg, err := config.PrepareAPIConfig()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	logger, err := logging.NewLogger(logging.WithConfig(cfg.LoggerConfig))
	if err != nil {
		panic(err)
	}

	err = tracer.InitTracer(cfg.JaegerConfig)
	if err != nil {
		logger.WithError(err).Error(ctx, "init tracer error")
	}

	store := storage.NewStorage()

	httpSrv, err := httpserver.NewHTTPServer(ctx, cfg, store, logger)
	if err != nil {
		logger.WithError(err).Fatal(ctx, "start of HTTP server failed")
	}

	// http server init
	go func() {
		logger.Info(ctx, "HTTP server starting...")
		err = httpSrv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Fatal(ctx, "failed listen and serve")
		}
	}()

	pprofSrv, err := httpserver.NewHTTPPPprofServer(cfg, logger)
	if err != nil {
		logger.WithError(err).Fatal(ctx, "failed to init instances")
	}

	go func() {
		logger.Info(ctx, "HTTP pprof server starting...")
		err = pprofSrv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Fatal(ctx, "failed to start http listening")
		}
	}()

	grpcListenConn, err := net.Listen("tcp", cfg.GRPCListenAddress)
	if err != nil {
		logger.WithError(err).WithField("listen", cfg.GRPCListenAddress).Fatal(ctx, err)
	}

	grpcSrv, err := grpcserver.NewServer(grpcListenConn, store, logger)
	if err != nil {
		logger.WithError(err).Fatal(ctx, err)
	}

	go func() {
		logger.Info(ctx, "GRPC server starting...")
		err = grpcSrv.ListenAndServe()
		if err != nil {
			logger.WithError(err).Fatal(ctx, err)
		}
	}()

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	logger.Warn(ctx, "shutdown application")

	grpcSrv.Close(ctx)

	if pprofSrv != nil {
		err = pprofSrv.Close(ctx)
		if err != nil {
			logger.WithError(err).Error(ctx, "stop HTTP pprof server failed")
		}
	}

	err = httpSrv.Close(ctx)
	if err != nil {
		logger.WithError(err).Error(ctx, "stop HTTP server failed")
	}

	log.Print("stopped")
}
