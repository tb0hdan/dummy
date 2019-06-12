package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/akhripko/dummy/healthcheck"
	"github.com/akhripko/dummy/log"
	"github.com/akhripko/dummy/metrics"
	"github.com/akhripko/dummy/options"
	"github.com/akhripko/dummy/prometheus"
	"github.com/akhripko/dummy/service"
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
	// register metrics
	metrics.Register()

	// prepare main context
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)
	var wg = &sync.WaitGroup{}

	// build main service
	srv := service.New(config.Port)
	// build prometheus service
	prometheusSrv := prometheus.New(config.PrometheusPort)
	// build healthcheck service
	healthChecks := []func() error{srv.HealthCheck, prometheusSrv.StateCheck}
	readinessChecks := []func() error{srv.ReadinessCheck, prometheusSrv.StateCheck}
	healthSrv := healthcheck.New(config.HealthCheckPort, healthChecks, readinessChecks)

	// run service
	healthSrv.Run(ctx, wg)
	prometheusSrv.Run(ctx, wg)
	srv.Run(ctx, wg)

	// wait while services work
	wg.Wait()
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
