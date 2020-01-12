# ex: make deploy [env=dev]
#
.PHONY: build clean deploy

env?=dev

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/foo foo/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/bar bar/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy -s $(env) --verbose
