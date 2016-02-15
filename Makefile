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

dbuilder:
	# golang 1.5 is required but available from jessie backports but backports
	# repo has low priority by default so we mount preinstall directory to
	# raise backports priority
	docker run --rm -it \
		-v `pwd`/debian/scripts:/dbuilder/preinstall.d:ro \
		-v `pwd`:/dbuilder/sources \
		seznam/dbuilder:debian_jessie-backports
