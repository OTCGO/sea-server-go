test:
	go test ./...

build:
	script/build.sh

install:
	script/install.sh

fmt:
	@gofmt -w .

.PHONY: clean gotool ca help build