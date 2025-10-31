PASSWORD_GENERATOR_SCRIPT = ./rabbit-mq/generate-password-file.sh

ADMIN_PASSWORD_FILE = ./rabbit-mq/secrets/admin-password
ADMIN_PASSWORD_HASH_FILE = ./rabbit-mq/secrets/admin-password-hash

generate-random-password:
	@openssl rand -base64 32 | tr -d "=+/" | cut -c1-32

generate-rabbitmq-secrets:
	@random_password=$$(make -s generate-random-password); \
	echo $$random_password > $(ADMIN_PASSWORD_FILE); \
	$(PASSWORD_GENERATOR_SCRIPT) $$random_password $(ADMIN_PASSWORD_HASH_FILE); \



