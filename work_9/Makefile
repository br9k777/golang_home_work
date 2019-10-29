APP?=work_9
RELEASE?=1.0.1
GOOS?=linux

COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# build:
# 	go build

test_prepare:
	go test --run=^TestWork9Prepare$

# Build

# build:
# 	test: build
# 	GL_TEST_RUN=1 ./golangci-lint run -v
# 	GL_TEST_RUN=1 ./golangci-lint run --fast --no-config -v --skip-dirs 'test/testdata_etc,pkg/golinters/goanalysis/(checker|passes)'
# 	GL_TEST_RUN=1 ./golangci-lint run --no-config -v --skip-dirs 'test/testdata_etc,pkg/golinters/goanalysis/(checker|passes)'
# 	GL_TEST_RUN=1 go test -v ./...

# .PHONY: test

# .PHONY: check
# check: prepare_golangci
# 	golangci-lint --vendor ./...

.PHONY: build
build: clean
	CGO_ENABLED=0 GOOS=${GOOS} go build \
		-ldflags "-X main.version=${RELEASE} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}" \
		-o bin/${GOOS}/${APP} 

.PHONY: clean
clean:
	@rm -f bin/${GOOS}/${APP}

# .PHONY: vendor
# vendor: prepare_dep
# 	dep ensure

# HAS_DEP := $(shell command -v dep;)
# HAS_METALINTER := $(shell command -v gometalinter;)

.PHONY: prepare_progress_bar
prepare_progress_bar:
ifndef HAS_DEP
	go get -u -v -d github.com/schollz/progressbar && \
	go install -v github.com/schollz/progressbar
endif

.PHONY: prepare_loger
prepare_loger:
ifndef HAS_DEP
	go get -u -v -d go.uber.org/zap && \
	go install -v go.uber.org/zap
endif

.PHONY: get_pakages
get_pakages: prepare_loger prepare_progress_bar

.PHONY: prepare_golangci
prepare_golangci:
ifndef HAS_METALINTER
	go get -u -v -d github.com/alecthomas/gometalinter && \
	go install -v github.com/alecthomas/gometalinter && \
	golangci-lint --install --update
endif

all: test_prepare prepare_progress_bar  prepare_loger