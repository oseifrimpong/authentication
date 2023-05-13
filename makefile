TAG := keyspecs/auth:latest

package:
	@echo "Building Auth Service docker image"
	@docker build -f deploy/docker/Dockerfile -t $(TAG) .

start:
	@echo "Starting Auth Service..."
	@echo "Generating Swagger Files"
	swag init
	@sh ./deploy/scripts/up.sh

stop:
	@echo "Stopping Auth Service..."
	@sh ./deploy/scripts/down.sh
	
.PHONY: build start stop
