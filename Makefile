FLAGS=-ldflags "-linkmode external -extldflags -static"
RUN=docker run --rm -v /dev/log:/dev/log -v $(shell pwd)/build:/build -v $(shell pwd):/go/src/gitlab.bridge.systems/ortb/reporter golang:1.9-alpine
BUILD=$(RUN) go build -ldflags "-X main.version=$(VERSION)"
REGISTRY=registry.bridge.systems:5005/ortb/reporter
VERSION=$(shell git describe --tags)
NPROC=1

.PHONY: build

clean:
	$(RUN) rm -rf ./build

test:
	$(RUN) go test gitlab.bridge.systems/ortb/reporter/...

build:
	$(BUILD) -o /build/reporter gitlab.bridge.systems/ortb/reporter

image:
	docker build -t $(REGISTRY):latest .

lint:
	docker run --rm -v /dev/log:/dev/log -v $(shell pwd)/build:/build -v $(shell pwd):/go/src/gitlab.bridge.systems/ortb/reporter dsociative/gobuilder gometalinter --deadline 15m --vendored-linters --vendor -j $(NPROC) --disable-all -E unused -E gosimple -E deadcode src/gitlab.bridge.systems/ortb/reporter/...

release:
	docker tag $(REGISTRY):latest $(REGISTRY):$(VERSION)
	docker push $(REGISTRY):$(VERSION)
	docker push $(REGISTRY):latest
