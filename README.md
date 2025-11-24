# IoT Distributed Architecture

A distributed IoT data collection and monitoring system built with Go, RabbitMQ, TimescaleDB, Prometheus, and Grafana. This project demonstrates a scalable architecture for collecting, processing, and visualizing sensor data and system metrics from multiple IoT devices.

## Architecture

The system follows a microservices architecture with the following data flow:

```
IoT Devices → MQTT → RabbitMQ → Workers → TimescaleDB/Prometheus → Grafana
```

### Components

- **Client**: Go application that simulates IoT devices, collecting sensor data (temperature, humidity) and system metrics (CPU, memory, disk, network usage)
- **RabbitMQ**: Message broker that receives MQTT messages and routes them to dedicated queues
- **Data Worker**: Consumes sensor data from RabbitMQ and stores it in TimescaleDB
- **Metrics Worker**: Consumes system metrics from RabbitMQ and forwards them to Prometheus
- **TimescaleDB**: Time-series database optimized for storing sensor data
- **Prometheus**: Metrics collection and storage system
- **Grafana**: Visualization platform with pre-configured dashboards for sensor data and system metrics

## Prerequisites

- Docker and Docker Compose: [Get Started](https://www.docker.com/get-started/)
- Go 1.24+ (for building and running clients locally): [Get Started](https://go.dev/doc/install)
- Make: [Download](https://www.gnu.org/software/make/)
- OpenSSL (for secret generation)

## Quick Start

### 1. Clone Project
```bash
git clone https://github.com/RicardoCenci/iot-distributed-architecture.git
```

### 2. Setup Project

Generate all required secrets and configuration files:

```bash
make setup-project
```

This command will:
- Generate random passwords for RabbitMQ users, Grafana admin, and TimescaleDB
- Create `.env` file from `.env.example`
- Create `client/config.toml` from `client/config.toml.example` with the generated client password

### 3. Start Services

Start all services using Docker Compose:

```bash
docker compose up -d
```

This will start:
- RabbitMQ (ports 5672, 15672, 15692, 1883)
- TimescaleDB (port 5432)
- Prometheus (port 9090)
- Grafana (port 3000)
- Data Worker
- Metrics Worker

### 3. Access Services

- **Grafana**: http://localhost:3000
  - Username: `admin`
  - Password: Check with `make get-grafana-admin-password`
- **RabbitMQ Management UI**: http://localhost:15672
  - Username: `admin`
  - Password: Check `./rabbit-mq/secrets/admin-password`
- **Prometheus**: http://localhost:9090

### 4. Run IoT Client

Run a single simulated device:

```bash
cd client
go run cmd/single_random/main.go
```

Or run multiple devices:

```bash
cd client
go run cmd/multiple_random/main.go -num-devices 10
```

## Configuration

### Client Configuration

The client configuration is in `client/config.toml`. Key settings:

- `[device]`: Device ID
- `[mqtt]`: MQTT broker connection settings
- `[mqtt.topics.data_json]`: Sensor data topic configuration
- `[mqtt.topics.metrics]`: System metrics topic configuration


## Services and Ports

| Service | Port | Description |
|---------|------|-------------|
| RabbitMQ AMQP | 5672 | AMQP protocol port |
| RabbitMQ Management | 15672 | Web management UI |
| RabbitMQ Prometheus | 15692 | Prometheus metrics endpoint |
| RabbitMQ MQTT | 1883 | MQTT protocol port |
| TimescaleDB | 5432 | PostgreSQL port |
| Prometheus | 9090 | Prometheus web UI |
| Grafana | 3000 | Grafana web UI |

## Data Flow

### Sensor Data Flow

1. IoT device collects sensor data (temperature, humidity)
2. Data is serialized using Protocol Buffers and base64 encoded
3. Published to RabbitMQ via MQTT on topic `iot.device.data.binary`
4. RabbitMQ routes message to `data-queue`
5. Data worker consumes from queue, deserializes, and stores in TimescaleDB
6. Grafana queries TimescaleDB to visualize sensor data

### Metrics Data Flow

1. IoT device collects system metrics (CPU, memory, disk, network)
2. Metrics are serialized using Protocol Buffers and base64 encoded
3. Published to RabbitMQ via MQTT on topic `iot.device.metrics`
4. RabbitMQ routes message to `metrics-queue`
5. Metrics worker consumes from queue, deserializes, and sends to Prometheus
6. Grafana queries Prometheus to visualize system metrics

## Makefile Commands

- `make setup-project`: Generate all secrets and configuration files
- `make generate-secrets`: Generate all password secrets
- `make generate-rabbitmq-secrets`: Generate RabbitMQ user passwords
- `make generate-grafana-secrets`: Generate Grafana admin password
- `make generate-data-worker-secrets`: Generate TimescaleDB password
- `make get-grafana-admin-password`: Display Grafana admin password

## Protocol Buffers

The project uses Protocol Buffers for efficient data serialization:

- **SensorData**: Contains sensor_id, humidity, temperature, and timestamp
- **MetricsData**: Contains sensor_id, cpu_usage, memory_usage, disk_usage, network_usage, and timestamp

To regenerate Go code from `.proto` files:

```bash
cd shared/proto
protoc --go_out=. --go_opt=paths=source_relative *.proto
```

## Monitoring

### Grafana Dashboards

Two pre-configured dashboards are available:

1. **Sensor Data Dashboard**: Visualizes temperature and humidity over time
2. **System Metrics Dashboard**: Visualizes CPU, memory, disk, and network usage

### Prometheus Metrics

Prometheus collects:
- RabbitMQ metrics (from RabbitMQ Prometheus endpoint)
- System metrics from IoT devices (via metrics worker)

