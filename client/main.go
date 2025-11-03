package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/RicardoCenci/iot-distributed-architecture/client/app"
	"github.com/RicardoCenci/iot-distributed-architecture/client/config"
	"github.com/RicardoCenci/iot-distributed-architecture/client/device"
	"github.com/RicardoCenci/iot-distributed-architecture/client/drivers"
	"github.com/RicardoCenci/iot-distributed-architecture/shared/logger"
)

const DEFAULT_CONFIG_FILE = "config.toml"

func main() {

	configFile := flag.String("config", DEFAULT_CONFIG_FILE, "config file")

	flag.Parse()

	config := config.NewConfig()

	if err := config.LoadFromTomlFile(*configFile); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := config.Validate(); err != nil {
		log.Fatalf("Failed to validate config: %v", err)
	}

	loggerConfig := logger.Config{
		Level: config.Log.Level,
		Source: logger.SourceConfig{
			Enabled:  config.Log.Source.Enabled,
			Relative: config.Log.Source.Relative,
			AsJSON:   config.Log.Source.AsJSON,
		},
	}
	logger := logger.NewSlogLogger(loggerConfig)

	device := device.NewDevice(
		config.Device.ID,
		drivers.NewRandomDataDriver(),
	)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	// Handle shutdown signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sig
		logger.Info("Received shutdown signal, shutting down...")
		cancel()
	}()

	app := app.NewApp(
		config,
		device,
		logger,
	)

	app.Run(ctx)

	logger.Info("Application shutdown complete")
	os.Exit(0)
}
