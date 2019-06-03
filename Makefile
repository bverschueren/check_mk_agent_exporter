BUILD_IMAGE ?= golang:1.12
BINARY = check_mk_exporter

VERSION = $(shell git describe --tags --always --dirty)
OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
TAG = $(VERSION)_$(OS)_$(ARCH)

BUILD_DIRS := bin bin/$(OS)_$(ARCH)
REGISTRY_PREFIX ?= bverschueren

container-build: bin/$(OS)_$(ARCH)/$(BINARY)
bin/$(OS)_$(ARCH)/$(BINARY): $(BUILD_DIRS)
	@sed                             \
	    -e 's|{ARCH}|$(ARCH)|g'      \
	    -e 's|{OS}|$(OS)|g'          \
		Dockerfile.in > .dockerfile-$(OS)_$(ARCH)
	@docker run			\
	--rm				\
	-v $$(pwd):/src			\
	-w /src				\
	-v $(pwd)/build:/go/bin/	\
	$(BUILD_IMAGE)			\
	go build -o bin/$(OS)_$(ARCH)/$(BINARY) .

container-image: container-build
	@docker build -t $(REGISTRY_PREFIX)/$(BINARY):$(TAG) -f .dockerfile-$(OS)_$(ARCH) .

container-clean:
	@docker rmi $(REGISTRY_PREFIX)/$(BINARY):$(TAG)

container-test: container-image
	$(eval FAKE := $(shell mktemp))
	$(eval CONTAINER_ID := $(shell docker run -d -p2112:2112 -v $(FAKE):/etc/check_mk_exporter/ssh.yaml $(REGISTRY_PREFIX)/$(BINARY):$(TAG)))
	@curl localhost:2112
	@docker stop $(CONTAINER_ID)

build: $(BUILD_DIRS)
	go build -v -o bin/$(OS)_$(ARCH)/$(BINARY) .

clean:
	rm -rf bin/

$(BUILD_DIRS):
	@mkdir -p $@
