BUILD_IMAGE ?= golang:1.12
BINARY = check_mk_exporter

VERSION = $(shell git describe --tags --always --dirty)
OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
TAG = $(VERSION)_$(OS)_$(ARCH)
UID := $(shell id -u)
GID := $(shell id -g)

BUILD_DIRS := bin
REGISTRY_PREFIX ?= bverschueren

builder-image:
	docker build -t $(BINARY)-builder .

container-image: $(BUILD_DIRS) builder-image
	@s2i build . $(BINARY)-builder $(BINARY)-intermediate

container-image-slim: binary
	@docker build -f Dockerfile-runtime -t $(REGISTRY_PREFIX)/$(BINARY):$(TAG) .

container-clean:
	@docker rmi -f $(REGISTRY_PREFIX)/$(BINARY):$(TAG)
	@docker rmi -f $(BINARY)-intermediate
	@docker rmi -f $(BINARY)-builder

container-test: container-image
	$(eval FAKE := $(shell mktemp -d))
	$(shell touch $(FAKE)/ssh.yaml)
	$(shell chmod -R 755 $(FAKE))
	$(eval CONTAINER_ID := $(shell docker run -d -p2112:2112 -v $(FAKE):/etc/check_mk_exporter/ $(BINARY)-intermediate))
	@curl localhost:2112
	@docker stop $(CONTAINER_ID)

dev-environment:
	mkdir -p ./docker/ssh/{client,server}
	chmod 700 ./docker/ssh/{client,server}
	yes y|ssh-keygen -t rsa -b 2038 -f ./docker/ssh/client/id_rsa -C dev-key -N ""
	cp ./docker/ssh/client/id_rsa.pub ./docker/ssh/server/authorized_keys

dev-environment-clean:
	rm -rf ./docker/ssh/*/*
	@docker-compose down

bin/$(BINARY): $(BUILD_DIRS) container-image
	@docker run $(BINARY)-intermediate /usr/libexec/s2i/save-artifacts | tar -xvf - -C bin/

binary: bin/$(BINARY)

clean:
	rm -rf bin/

$(BUILD_DIRS):
	mkdir -p $@
