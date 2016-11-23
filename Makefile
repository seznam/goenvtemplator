all: build

TAG :=-$(shell git describe --tags)
ifeq "$(TAG)" "-"
TAG :=
endif

LDFLAGS :=-X main.buildVersion=$(TAG)

.PHONY:=all build test release clean


build:
	go build -ldflags "$(LDFLAGS)"

test:
	go test

release:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o goenvtemplator-amd64 
	tar -cJf goenvtemplator.tar.xz goenvtemplator-amd64
	tar -tvf goenvtemplator.tar.xz

clean:
	$(RM) goenvtemplator{,-amd64,.tar.xz}
