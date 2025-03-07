.PHONY: vet test localtest run local_env_run migrate get up down

# docker version
version := 0.1.0
containername=slyip
dbname=slyip
schema=slyip

SERVICE_BUILD_PATH=tmp/main

# Read environment
env_local=local
env=$(env_local)

ifneq (,$(wildcard ./.$(env).env))
    include .$(env).env
    export
endif


# go
get:
	go get -d -v ./...
vet:
	go mod tidy
	go vet ./...
	go fmt ./...

run:
	rm -f $(SERVICE_BUILD_PATH)
	CGO_ENABLED=0 go build -o $(SERVICE_BUILD_PATH) cmd/service/main.go
	chmod +x $(SERVICE_BUILD_PATH)
	$(SERVICE_BUILD_PATH)
# Docker
build_docker:
	docker buildx build --platform="linux/amd64"  -f Dockerfile -t  $(containername):$(version) .
push_docker: build_docker
	docker tag $(containername):$(version) leondroid/$(containername):$(version)
	docker push leondroid/$(containername):$(version)
# Database
up:
	cd internal/goose/migrations; goose  -v postgres "postgres://$(postgres_user):$(postgres_pwd)@$(postgres_url)/$(dbname)?sslmode=disable" up
down:
	cd internal/goose/migrations; goose  postgres "postgres://$(postgres_user):$(postgres_pwd)@$(postgres_url)/$(dbname)?sslmode=disable" down
jet:
	jet -dsn=postgres://$(postgres_user):$(postgres_pwd)@$(postgres_url)/$(dbname)?sslmode=disable -schema=$(schema) -path=./.gen

recreate: down up jet
rerun: recreate run
