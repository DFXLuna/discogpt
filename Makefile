lint:
	golangci-lint run ./...

build:
	docker build . -f ./deploy/Dockerfile -t ghcr.io/dfxluna/discogpt:latest

push:
	docker push ghcr.io/dfxluna/discogpt:latest