package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

type SensorData struct {
	DeviceID    string
	Timestamp   time.Time
	Humidity    float32
	Temperature float32
}

func NewDatabase(connectionString string) (*Database, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{db: db}

	if err := database.runMigrations(connectionString); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return database, nil
}

func (d *Database) runMigrations(connectionString string) error {
	migrationURL, err := convertToMigrationURL(connectionString)
	if err != nil {
		return fmt.Errorf("failed to convert connection string: %w", err)
	}

	m, err := migrate.New(
		"file://migrations",
		migrationURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func convertToMigrationURL(connectionString string) (string, error) {
	parts := make(map[string]string)
	pairs := strings.Fields(connectionString)

	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			parts[kv[0]] = kv[1]
		}
	}

	host := parts["host"]
	port := parts["port"]
	user := parts["user"]
	password := parts["password"]
	dbname := parts["dbname"]
	sslmode := parts["sslmode"]

	if host == "" || user == "" || dbname == "" {
		return "", fmt.Errorf("missing required connection parameters")
	}

	if port == "" {
		port = "5432"
	}

	if sslmode == "" {
		sslmode = "disable"
	}

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   dbname,
	}

	q := u.Query()
	q.Set("sslmode", sslmode)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (d *Database) InsertSensorData(data SensorData) error {
	query := `INSERT INTO sensor_data (time, device_id, humidity, temperature) VALUES ($1, $2, $3, $4)`
	_, err := d.db.Exec(query, data.Timestamp, data.DeviceID, data.Humidity, data.Temperature)
	if err != nil {
		return fmt.Errorf("failed to insert sensor data: %w", err)
	}
	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
