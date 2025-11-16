package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/RicardoCenci/iot-distributed-architecture/shared/logger"
)

type Config struct {
	User      string `json:"rabbitmq_user"`
	Password  string
	Domain    string        `json:"rabbitmq_domain"`
	Port      string        `json:"rabbitmq_port"`
	QueueName string        `json:"rabbitmq_queue_name"`
	Log       logger.Config `json:"log"`
}

var DEFAULT_SECRET_PATH = getStringEnv("DEFAULT_SECRET_PATH", "/run/secrets/")

func NewConfig() *Config {
	fileConfig, err := getFromFile("config.json")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
		return nil
	}

	password, err := getFromSecret("RABBITMQ_DATA_WORKER_PASSWORD")

	if err != nil {
		log.Fatalf("Failed to read secret RABBITMQ_DATA_WORKER_PASSWORD: %v", err)
	}

	return &Config{
		User:      getStringEnv("RABBITMQ_DATA_WORKER_USER", fileConfig.User),
		Password:  password,
		Domain:    getStringEnv("RABBITMQ_DOMAIN", fileConfig.Domain),
		Port:      getStringEnv("RABBITMQ_AMQP_PORT", fileConfig.Port),
		QueueName: getStringEnv("RABBITMQ_DATA_WORKER_QUEUE_NAME", fileConfig.QueueName),
		Log: logger.Config{
			Level: getStringEnv("LOG_LEVEL", "debug"),
			Source: logger.SourceConfig{
				Enabled:  getBoolEnv("LOG_SOURCE_ENABLED", true),
				Relative: getBoolEnv("LOG_SOURCE_RELATIVE", true),
				AsJSON:   getBoolEnv("LOG_SOURCE_AS_JSON", false),
			},
		},
	}
}

func getFromFile(path string) (*Config, error) {
	config, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var configData Config
	json.Unmarshal(config, &configData)

	return &configData, nil
}

func getBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true"
}

func getStringEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getFromSecret(name string) (string, error) {
	path := filepath.Join(DEFAULT_SECRET_PATH, name)

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
