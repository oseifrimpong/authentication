#!/bin/sh
# Charles Cyril Nettey <cyril@keyspecs.com>
# This is the third and final script of the entryPoint call sequence

echo "auth-api entrypoint"
set -e

# Call command issued to the docker service
echo "auth-service exec: $@"
exec "$@"