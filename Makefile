all: build

buildprod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/web/web ./cmd/web
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/api/api ./cmd/api
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/crawler/crawler ./cmd/crawler

build:
	go build -o ./cmd/web/web ./cmd/web
	go build -o ./cmd/api/api ./cmd/api
	go build -o ./cmd/crawler/crawler ./cmd/crawler
