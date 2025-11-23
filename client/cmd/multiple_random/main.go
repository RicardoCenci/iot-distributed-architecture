package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
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
	numDevices := flag.Int("num-devices", 10, "number of devices to run")

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

	logger.Info("Starting multiple random devices", "num-devices", *numDevices)

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

	wg := sync.WaitGroup{}
	wg.Add(*numDevices)

	for i := 0; i < *numDevices; i++ {
		go func(i int) {
			defer wg.Done()

			deviceLogger := logger.WithContext(
				"device-id", fmt.Sprintf("device-%d", i),
			)

			deviceLogger.Info("Starting device")

			device := device.NewDevice(
				fmt.Sprintf("device-%d", i),
				drivers.NewRandomDataDriver(),
			)

			app := app.NewApp(
				config,
				device,
				deviceLogger,
			)

			app.Run(ctx)
		}(i)
	}
	wg.Wait()

	logger.Info("Application shutdown complete")
	os.Exit(0)
}
