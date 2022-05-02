package main

import (
	"cow/internal"
	http_internal "cow/internal/http"
	"fmt"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	//logger, _ := initLogger()
	//tntClient := initTntClient(logger)
	//scoreStorage := tnt.NewScoresStorage(tntClient)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		//initHttpServer(scoreStorage, logger)
		//wg.Done()
	}()

	wg.Wait()
}

// initLogger inits and returns logger
func initLogger() (*zap.Logger, error) {
	lvl := zap.InfoLevel
	err := lvl.UnmarshalText([]byte(os.Getenv("CL_LOG_LEVEL")))
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal log-level: %w", err)
	}
	opts := zap.NewProductionConfig()
	opts.Level = zap.NewAtomicLevelAt(lvl)
	opts.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if os.Getenv("CL_LOG_JSON") == "" {
		opts.Encoding = "console"
		opts.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return opts.Build()
}

// initTntClient inits and returns tarantool client
func initTntClient(logger *zap.Logger) *tarantool.Connection {
	var tarantoolClient *tarantool.Connection
	tarantoolOpts := tarantool.Opts{
		Timeout:   5 * time.Second,
		Reconnect: 1 * time.Second,
		User:      os.Getenv("CL_TARANTOOL_USERNAME"),
		Pass:      os.Getenv("CL_TARANTOOL_PASSWORD"),
	}

	logger.Info(
		"tnt",
		zap.String("username", os.Getenv("CL_TARANTOOL_USERNAME")),
		zap.String("password", os.Getenv("CL_TARANTOOL_PASSWORD")),
		zap.String("host", os.Getenv("CL_TARANTOOL_HOST")),
		zap.String("port", os.Getenv("CL_TARANTOOL_PORT")),
	)

	tarantoolClient, err := tarantool.Connect(
		os.Getenv("CL_TARANTOOL_HOST")+":"+os.Getenv("CL_TARANTOOL_PORT"),
		tarantoolOpts,
	)
	if err != nil {
		logger.Panic("failed connect to tarantool", zap.Error(err))
		return nil
	}

	logger.Info("tarantool connection established")

	return tarantoolClient
}

// initHttpServer inits http server
func initHttpServer(scoreStorage internal.ScoresStorage, logger *zap.Logger) {
	resp := http_internal.NewResponse(logger)
	api := http_internal.NewApi(scoreStorage, logger, resp)

	logger.Info(
		"Service HTTP-server is starting",
		zap.String("port", os.Getenv("CL_HTTP_SERVICE_LISTEN")),
	)

	http.HandleFunc("/v1/ping", api.Ping)
	http.HandleFunc("/v1/add", api.SetScore)
	http.HandleFunc("/v1/get", api.FindScore)
	err := http.ListenAndServe(os.Getenv("CL_HTTP_SERVICE_LISTEN"), nil)
	if err != nil {
		logger.Fatal("can't listen and serve", zap.Error(err))
		return
	}
}
