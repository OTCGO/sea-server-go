test:
	go test ./...

build:
	script/build.sh

fmt:
	@gofmt -w .

.PHONY: clean gotool ca help build