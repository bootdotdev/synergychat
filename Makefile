all: build

deploy: buildprod builddocker pushdocker

buildprod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/web/web ./cmd/web
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/api/api ./cmd/api
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/crawler/crawler ./cmd/crawler

build:
	go build -o ./cmd/web/web ./cmd/web
	go build -o ./cmd/api/api ./cmd/api
	go build -o ./cmd/crawler/crawler ./cmd/crawler

builddocker:
	docker build -t lanecwagner/synergychat-web ./cmd/web
	docker build -t lanecwagner/synergychat-api ./cmd/api
	docker build -t lanecwagner/synergychat-crawler ./cmd/crawler

pushdocker:
	docker push lanecwagner/synergychat-web
	docker push lanecwagner/synergychat-api
	docker push lanecwagner/synergychat-crawler
