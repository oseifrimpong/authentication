#!/bin/sh
# Charles Cyril Nettey <cyril@keyspecs.com>
# This script expects an env var passed in the docker runtime used to query secrets from AWS Secrets Manager

set -e

echo "Preparing Auth Keys"
if [[  "$SKIP_AUTH_KEY_FETCH" == "1" ]]; then
  echo "SKIP_AUTH_KEY_FETCH env is set, assuming key files exist..."
  exec "$@"
fi

if [[ -z "$AUTH_KEYS_BUCKET" || -z $AUTH_KEYS_PATH ]]; then
    echo "AUTH_KEYS_BUCKET: $AUTH_KEYS_BUCKET"
    echo "AUTH_KEYS_PATH: $AUTH_KEYS_PATH"
    echo "MISSING CREDENTIALS, Exiting..."
    exit 1
fi

if [[ -z "$CUSTOM_S3_ENDPOINT" ]];then
  KEY_FILES=$(aws s3api list-objects --bucket $AUTH_KEYS_BUCKET | jq -r '.Contents[].Key')
  # For each item in the KEY_FILES API Content list, get object from aws s3 api
  for KEY_FILE in $KEY_FILES; do
      echo "Downloading $KEY_FILE..."
      aws s3 cp s3://$AUTH_KEYS_BUCKET/$KEY_FILE $AUTH_KEYS_PATH/$KEY_FILE
  done

else

  KEY_FILES=$(aws s3api --endpoint-url $CUSTOM_S3_ENDPOINT list-objects --bucket $AUTH_KEYS_BUCKET | jq -r '.Contents[].Key')
  # For each item in the KEY_FILES API Content list, get object from aws s3 api
  for KEY_FILE in $KEY_FILES; do
      echo "Downloading $KEY_FILE..."
      aws s3 --endpoint-url $CUSTOM_S3_ENDPOINT cp s3://$AUTH_KEYS_BUCKET/$KEY_FILE $AUTH_KEYS_PATH/$KEY_FILE
  done

fi


echo "Done."

echo "Starting $@..."
exec "$@"