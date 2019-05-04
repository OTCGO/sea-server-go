test:
	go test ./...

build:
	script/build.sh

install:
	script/install.sh

fmt:
	@gofmt -w .

image:
	docker build -t sea-server-go .

.PHONY: clean gotool ca help build