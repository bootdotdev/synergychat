all: build

buildprod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/web ./cmd/web
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api ./cmd/api
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/crawler ./cmd/crawler

build:
	go build -o bin/web ./cmd/web
	go build -o bin/api ./cmd/api
	go build -o bin/crawler ./cmd/crawler
