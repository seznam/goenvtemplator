all: build

TAG :=-$(shell git describe --tags)
ifeq "$(TAG)" "-"
TAG :=
endif
DESTDIR      ?=
BINARY       ?= goenvtemplator2

LDFLAGS :=-X main.buildVersion=$(TAG)

.PHONY: all build test release clean install lint

build:
	go build -o $(BINARY) -ldflags "$(LDFLAGS)"

install:
	install -D -m 0755 $(BINARY) $(DESTDIR)/usr/bin/$(BINARY)

test:
	go test

lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.42.0 golangci-lint run -v

release:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o goenvtemplator2-amd64
	tar -cJf goenvtemplator.tar.xz goenvtemplator2-amd64
	tar -tvf goenvtemplator.tar.xz

clean:
	$(RM) goenvtemplator{2,2-amd64,.tar.xz}
