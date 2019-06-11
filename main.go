package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/akhripko/dummy/log"
	"github.com/akhripko/dummy/metrics"
	"github.com/akhripko/dummy/options"
	"github.com/akhripko/dummy/service"
	"github.com/akhripko/dummy/service_healthcheck"
	"github.com/akhripko/dummy/service_prometheus"
)

func main() {
	config := options.ReadEnv()
	err := initLogger(config)
	if err != nil {
		fmt.Println("log initialization error:", err.Error())
		os.Exit(1)
		return
	}

	// register metrics
	metrics.Register()

	// prepare context
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)
	var wg = &sync.WaitGroup{}

	// main service
	// TODO: fix it
	var srv service.Service

	healthChecks := []func() error{}
	readinessChecks := []func() error{}

	// health check service
	prometheusSrv := service_prometheus.New(config.PrometheusPort)
	prometheusSrv.Run(ctx, wg)

	// health check service
	healthSrv := service_healthcheck.New(config.HealthCheckPort, healthChecks, readinessChecks)
	healthSrv.Run(ctx, wg)

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
