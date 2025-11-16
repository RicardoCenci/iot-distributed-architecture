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

# Password for the Grafana admin to access the Grafana UI
GF_SECURITY_ADMIN_PASSWORD_FILE = ./grafana/secrets/admin-password

generate-random-password:
	@openssl rand -base64 32 | tr -d "=+/" | cut -c1-32

generate-rabbitmq-secrets:
	@mkdir -p ./rabbit-mq/secrets; \
	random_password=$$(make -s generate-random-password); \
	echo $$random_password > $(ADMIN_PASSWORD_FILE); \
	$(PASSWORD_GENERATOR_SCRIPT) $$random_password $(ADMIN_PASSWORD_HASH_FILE); \
	echo $$random_password > $(METRICS_PASSWORD_FILE); \
	$(PASSWORD_GENERATOR_SCRIPT) $$random_password $(METRICS_PASSWORD_HASH_FILE); \
	echo $$random_password > $(DATA_WORKER_PASSWORD_FILE); \
	$(PASSWORD_GENERATOR_SCRIPT) $$random_password $(DATA_WORKER_PASSWORD_HASH_FILE); \
	echo $$random_password > $(CLIENT_USER_PASSWORD_FILE); \
	$(PASSWORD_GENERATOR_SCRIPT) $$random_password $(CLIENT_USER_PASSWORD_HASH_FILE); \
	
generate-grafana-secrets:
	@mkdir -p ./grafana/secrets; \
	random_password=$$(make -s generate-random-password); \
	echo $$random_password > $(GF_SECURITY_ADMIN_PASSWORD_FILE); \

generate-secrets:
	@make generate-rabbitmq-secrets; \
	make generate-grafana-secrets; \
