package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/akhripko/dummy/cache"

	"github.com/akhripko/dummy/healthcheck"
	"github.com/akhripko/dummy/log"
	"github.com/akhripko/dummy/metrics"
	"github.com/akhripko/dummy/options"
	"github.com/akhripko/dummy/prometheus"
	"github.com/akhripko/dummy/service"
	"github.com/akhripko/dummy/storage"
)

func main() {
	// read service config from os env
	config := options.ReadEnv()
	// init logger
	err := initLogger(config)
	if err != nil {
		fmt.Println("log initialization error:", err.Error())
		os.Exit(1)
		return
	}
	log.Info("begin...")
	// register metrics
	metrics.Register()

	// prepare main context
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)
	var wg = &sync.WaitGroup{}

	// build db
	db, err := storage.NewSQLDB(ctx, storage.SQLDBConfig(config.SQLDB))
	if err != nil {
		log.Error("sql db init error:", err.Error())
		os.Exit(1)
	}
	// build cache
	cc, err := cache.NewRedis(ctx, config.CacheAddr)
	if err != nil {
		log.Error("cache init error:", err.Error())
		os.Exit(1)
	}

	// build main service
	srv := service.New(config.Port, db, cc)
	// build prometheus service
	prometheusSrv := prometheus.New(config.PrometheusPort)
	// build healthcheck service
	healthChecks := []func() error{srv.HealthCheck, prometheusSrv.HealthCheck}
	healthSrv := healthcheck.New(config.HealthCheckPort, healthChecks)

	// run service
	healthSrv.Run(ctx, wg)
	prometheusSrv.Run(ctx, wg)
	srv.Run(ctx, wg)

	// wait while services work
	wg.Wait()
	log.Info("end")
}

func initLogger(config *options.Config) error {
	logLvl := log.StrToLogLevel(config.LogLevel)
	logger, err := log.NewConsoleLogger(logLvl)
	if err != nil {
		return err
	}
	log.SetLogger(logger)
	return nil
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.Error("Got Interrupt signal")
		stop()
	}()
}
