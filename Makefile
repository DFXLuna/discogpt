lint:
	golangci-lint run ./...

build:
	docker build . -f ./deploy/Dockerfile -t dfxluna/discogpt:latest

push:
	docker push dfxluna/discogpt:latest