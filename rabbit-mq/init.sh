#!/bin/bash
set -e

for secret_file in /run/secrets/*; do
  if [ -f "$secret_file" ]; then
    secret_name=$(basename "$secret_file")
    secret_value=$(cat "$secret_file")
    export "$secret_name=$secret_value"
  fi
done

envsubst < /etc/rabbitmq/definitions.template.json > /etc/rabbitmq/definitions.json

rabbitmq-server
