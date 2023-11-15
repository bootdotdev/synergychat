all: build

deploy: buildprod builddocker pushdocker

buildprod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/web/web ./cmd/web
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/api/api ./cmd/api
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/crawler/crawler ./cmd/crawler
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/testcpu/testcpu ./cmd/testcpu
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/testram/testram ./cmd/testram

build:
	go build -o ./cmd/web/web ./cmd/web
	go build -o ./cmd/api/api ./cmd/api
	go build -o ./cmd/crawler/crawler ./cmd/crawler
	go build -o ./cmd/testcpu/testcpu ./cmd/testcpu
	go build -o ./cmd/testram/testram ./cmd/testram

builddocker:
	docker build -t lanecwagner/synergychat-web ./cmd/web
	docker build -t lanecwagner/synergychat-api ./cmd/api
	docker build -t lanecwagner/synergychat-crawler ./cmd/crawler
	docker build -t lanecwagner/synergychat-testram ./cmd/testram
	docker build -t lanecwagner/synergychat-testcpu ./cmd/testcpu

pushdocker:
	docker push lanecwagner/synergychat-web
	docker push lanecwagner/synergychat-api
	docker push lanecwagner/synergychat-crawler
	docker push lanecwagner/synergychat-testram
	docker push lanecwagner/synergychat-testcpu
