CREATE TABLE IF NOT EXISTS sensor_data (
	time TIMESTAMPTZ NOT NULL,
	device_id TEXT NOT NULL,
	humidity REAL NOT NULL,
	temperature REAL NOT NULL
);

