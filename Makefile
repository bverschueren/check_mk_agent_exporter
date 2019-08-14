BUILD_IMAGE ?= golang:1.12
BINARY = check_mk_exporter

VERSION = $(shell git describe --tags --always --dirty)
OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
TAG = $(VERSION)_$(OS)_$(ARCH)
UID := $(shell id -u)
GID := $(shell id -g)

BUILD_DIRS := bin bin/$(OS)_$(ARCH)
REGISTRY_PREFIX ?= bverschueren

container-build: bin/$(OS)_$(ARCH)/$(BINARY)
bin/$(OS)_$(ARCH)/$(BINARY): $(BUILD_DIRS)
	@docker run			\
	--rm				\
	-v $$(pwd):/src			\
	-w /src				\
	-v $(pwd)/build:/go/bin/	\
	$(BUILD_IMAGE)			\
	go build -o bin/$(OS)_$(ARCH)/$(BINARY) .

container-image: container-build
	@docker build -t $(REGISTRY_PREFIX)/$(BINARY):$(TAG) \
		--build-arg=OS=$(OS) \
		--build-arg=ARCH=$(ARCH) \
		--build-arg=UID=$(UID) \
		--build-arg=GID=$(GID) \
		.

container-clean:
	@docker rmi $(REGISTRY_PREFIX)/$(BINARY):$(TAG)

container-test: container-image
	$(eval FAKE := $(shell mktemp))
	$(eval CONTAINER_ID := $(shell docker run -d -p2112:2112 -v $(FAKE):/etc/check_mk_exporter/ssh.yaml $(REGISTRY_PREFIX)/$(BINARY):$(TAG)))
	@curl localhost:2112
	@docker stop $(CONTAINER_ID)

dev-environment:
	mkdir -p ./docker/ssh/{client,server}
	chmod 700 ./docker/ssh/{client,server}
	yes y|ssh-keygen -t rsa -b 2038 -f ./docker/ssh/client/id_rsa -C dev-key -N ""
	cp ./docker/ssh/client/id_rsa.pub ./docker/ssh/server/authorized_keys

clean-dev-environment:
	rm -rf ./docker/ssh/*/*
	@docker image ls|grep check_mk_exporter|awk '{print $$3}'|xargs -r docker image rm --force

build: $(BUILD_DIRS)
	go build -v -o bin/$(OS)_$(ARCH)/$(BINARY) .

clean:
	rm -rf bin/

$(BUILD_DIRS):
	mkdir -p $@
