test:
	go test ./...

build:
	script/build.sh

install:
	script/install.sh

fmt:
	@gofmt -w .

image:
	docker build -t kiang/sea-server-go .

.PHONY: clean gotool ca help build