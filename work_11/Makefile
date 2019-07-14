APP?=app
RELEASE?=0.0.1
GOOS?=linux

COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Build

# build:
# 	test: build
# 	GL_TEST_RUN=1 ./golangci-lint run -v
# 	GL_TEST_RUN=1 ./golangci-lint run --fast --no-config -v --skip-dirs 'test/testdata_etc,pkg/golinters/goanalysis/(checker|passes)'
# 	GL_TEST_RUN=1 ./golangci-lint run --no-config -v --skip-dirs 'test/testdata_etc,pkg/golinters/goanalysis/(checker|passes)'
# 	GL_TEST_RUN=1 go test -v ./...

# .PHONY: test

.PHONY: check
check: prepare_golangci
	golangci-lint --vendor ./...

.PHONY: build
build: clean
	CGO_ENABLED=0 GOOS=${GOOS} go build \
		-ldflags "-X main.version=${RELEASE} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}" \
		-o bin/${GOOS}/${APP} 

.PHONY: clean
clean:
	@rm -f bin/${GOOS}/${APP}

.PHONY: vendor
vendor: prepare_dep
	dep ensure

HAS_DEP := $(shell command -v dep;)
HAS_METALINTER := $(shell command -v gometalinter;)

.PHONY: prepare_viper
prepare_viper:
ifndef HAS_DEP
	go get -u -v -d github.com/spf13/viper && \
	go install -v github.com/spf13/viper
endif

.PHONY: prepare_golangci
prepare_golangci:
ifndef HAS_METALINTER
	go get -u -v -d github.com/alecthomas/gometalinter && \
	go install -v github.com/alecthomas/gometalinter && \
	golangci-lint --install --update
endif