PASSWORD_GENERATOR_SCRIPT = ./rabbit-mq/generate-password-file.sh

# Password for the admin to access the RabbitMQ
ADMIN_PASSWORD_FILE = ./rabbit-mq/secrets/admin-password
ADMIN_PASSWORD_HASH_FILE = ./rabbit-mq/secrets/admin-password-hash

# Password for the Prometheus user to access the RabbitMQ
METRICS_PASSWORD_FILE = ./rabbit-mq/secrets/metrics-password
METRICS_PASSWORD_HASH_FILE = ./rabbit-mq/secrets/metrics-password-hash

# Password for the data worker to access the RabbitMQ
DATA_WORKER_PASSWORD_FILE = ./rabbit-mq/secrets/data-worker-password
DATA_WORKER_PASSWORD_HASH_FILE = ./rabbit-mq/secrets/data-worker-password-hash

# Password for the client user to access the RabbitMQ
CLIENT_USER_PASSWORD_FILE = ./rabbit-mq/secrets/client-user-password
CLIENT_USER_PASSWORD_HASH_FILE = ./rabbit-mq/secrets/client-user-password-hash

# Password for the metrics worker to access the RabbitMQ
METRICS_WORKER_PASSWORD_FILE = ./rabbit-mq/secrets/metrics-worker-password
METRICS_WORKER_PASSWORD_HASH_FILE = ./rabbit-mq/secrets/metrics-worker-password-hash

# Password for the data worker to access the TimescaleDB
TIMESCALEDB_PASSWORD_FILE = ./workers/data/secrets/timescaledb-password
TIMESCALEDB_PASSWORD_HASH_FILE = ./workers/data/secrets/timescaledb-password-hash

# Password for the Grafana admin to access the Grafana UI
GF_SECURITY_ADMIN_PASSWORD_FILE = ./grafana/secrets/admin-password

generate-random-password:
	@openssl rand -base64 32 | tr -d "=+/" | cut -c1-32

generate-rabbitmq-secrets:
	@mkdir -p ./rabbit-mq/secrets; \
	random_password=$$(make -s generate-random-password); \
	printf '%s' $$random_password > $(ADMIN_PASSWORD_FILE); \
	$(PASSWORD_GENERATOR_SCRIPT) $$random_password $(ADMIN_PASSWORD_HASH_FILE); \
	printf '%s' $$random_password > $(METRICS_PASSWORD_FILE); \
	$(PASSWORD_GENERATOR_SCRIPT) $$random_password $(METRICS_PASSWORD_HASH_FILE); \
	printf '%s' $$random_password > $(DATA_WORKER_PASSWORD_FILE); \
	$(PASSWORD_GENERATOR_SCRIPT) $$random_password $(DATA_WORKER_PASSWORD_HASH_FILE); \
	printf '%s' $$random_password > $(CLIENT_USER_PASSWORD_FILE); \
	$(PASSWORD_GENERATOR_SCRIPT) $$random_password $(CLIENT_USER_PASSWORD_HASH_FILE); \
	printf '%s' $$random_password > $(METRICS_WORKER_PASSWORD_FILE); \
	$(PASSWORD_GENERATOR_SCRIPT) $$random_password $(METRICS_WORKER_PASSWORD_HASH_FILE); \
	
generate-grafana-secrets:
	@mkdir -p ./grafana/secrets; \
	random_password=$$(make -s generate-random-password); \
	printf '%s' $$random_password > $(GF_SECURITY_ADMIN_PASSWORD_FILE); \

generate-data-worker-secrets:
	@mkdir -p ./workers/data/secrets; \
	random_password=$$(make -s generate-random-password); \
	printf '%s' $$random_password > $(TIMESCALEDB_PASSWORD_FILE);

generate-secrets:
	@make -s generate-rabbitmq-secrets; \
	make -s generate-grafana-secrets; \
	make -s generate-data-worker-secrets; \
	echo "Secrets generated successfully"; \
	echo "Your Grafana admin username: admin"; \
	echo "Your Grafana admin password: $$(cat $(GF_SECURITY_ADMIN_PASSWORD_FILE))"; \


setup-project:
	@make -s generate-secrets; \
	cp .env.example .env; \
	cp client/config.toml.example client/config.toml; \
	sed -i "s|<YOUR_PASSWORD>|$$(cat $(CLIENT_USER_PASSWORD_FILE))|g" client/config.toml; \

get-grafana-admin-password:
	@echo "Your Grafana admin password: $$(cat $(GF_SECURITY_ADMIN_PASSWORD_FILE))"; \