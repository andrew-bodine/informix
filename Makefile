# TODO: Add some title information here.

# This Makefile encapsulates all the details of how to build, test, and deploy
# Informix across different environments.

# Constants
CONTAINER_IMAGE=informix
CONTAINER_INSTANCE=informix

# Default entrypoint (i.e. `$ make`)
all: build

# Create a container image for portability.
build:
	docker build --force-rm -t $(CONTAINER_IMAGE) .

# Create a cross-compiled runnable suitable for a Raspberry Pi and the arm
# architecture.
pi: build
	docker run --rm \
		-v $$(pwd):/root/build \
		-e GOARCH=arm \
		-e GOOS=linux \
		--entrypoint go \
		$(CONTAINER_IMAGE) \
		build -o /root/build/$(CONTAINER_INSTANCE) main.go

# Run a container from the built image, re-creating if necessary. We pull in
# environment variables for configuring each container, the invoker should set
# these (i.e. `$ export WIOT_ORG_ID=...`)
run: clean build
	docker run --name $(CONTAINER_INSTANCE) \
	 	-d \
		$(CONTAINER_IMAGE) \
		/go/bin/$(CONTAINER_IMAGE)
#		-e WIOT_ORG_ID=${WIOT_ORG_ID} \
#		-e WIOT_AUTH_TOKEN=${WIOT_AUTH_TOKEN} \
#		-e WIOT_DEVICE_ID=${WIOT_DEVICE_ID} \
#		-e WIOT_DEVICE_TYPE=${WIOT_DEVICE_TYPE} \

# Remove the possible container instance.
clean:
	docker rm -fv $(CONTAINER_INSTANCE) || true

# Remove the possible container instance, and container image.
clean-all: clean
	docker rmi $(CONTAINER_IMAGE) || true
