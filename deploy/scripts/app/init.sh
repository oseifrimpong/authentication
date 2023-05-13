#!/bin/sh
# Charles Cyril Nettey <cyril@keyspecs.com>
# This is the first script of the entrypoint call sequence
# It is responsible for building the .env file for the Go Service

set -a

if [[ "$SKIP_SECRETS_FETCH" == "1" ]]; then
  echo "SKIP_SECRETS_FETCH env is set, assuming env vars have been set..."
  exec "$@"
fi

if [[ -z "$SECRET_SERVICE" || -z "$AWS_REGION" ]]; then
  echo "SECRET_SERVICE is not set. Exitting..."
  exit 1
fi

echo "Reading Secrets from AWS Secrets Manager at $SECRET_SERVICE "
aws secretsmanager get-secret-value --region $AWS_REGION --secret-id $SECRET_SERVICE --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value|tostring)")|.[]' > .env

echo "Loading Secrets into OS Env"
source .env

echo "Testing Application ENV"
echo APP_PORT="$APP_PORT"
echo GIN_MODE="$GIN_MODE"

echo "Starting $@..."
exec "$@"