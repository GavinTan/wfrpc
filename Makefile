export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

all: build

build: wfrpc

wfrpc:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/wfrpc ./cmd/wfrpc

clean:
	rm -rf ./bin
