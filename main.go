package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/akhripko/dummy/log"
	"github.com/akhripko/dummy/options"
)

func main() {
	config := options.ReadEnv()
	err := initLogger(config)
	if err != nil {
		fmt.Println("log initialization error:", err.Error())
		os.Exit(1)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)
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
