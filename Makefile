
# Variables
APP_NAME=apex-load-generator
DOCKER_IMAGE=$(APP_NAME):latest
BUILD_FLAGS := "-X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash=`git rev-parse HEAD` -X main.version=${VERSION} -linkmode external"
GOOS = linux
GOARCH = amd64

# Build the Go application
build:
	go build -ldflags ${BUILD_FLAGS} -o out/${shell uname -s}/${APP_NAME} .

# Build the Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# Run the Docker container
docker-run:
	docker run --rm -it $(DOCKER_IMAGE)

# Clean the build artifacts
clean:
	rm -f $(APP_NAME)

# Phony targets
.PHONY: build docker-build docker-run clean