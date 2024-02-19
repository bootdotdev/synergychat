all: build

deploy: buildprod builddocker pushdocker

buildprod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/web/web-linux-amd64 ./cmd/web
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./cmd/web/web-linux-arm64 ./cmd/web
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/api/api-linux-amd64 ./cmd/api
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./cmd/api/api-linux-arm64 ./cmd/api
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/crawler/crawler-linux-amd64 ./cmd/crawler
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./cmd/crawler/crawler-linux-arm64 ./cmd/crawler
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/testcpu/testcpu-linux-amd64 ./cmd/testcpu
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./cmd/testcpu/testcpu-linux-arm64 ./cmd/testcpu
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/testram/testram-linux-amd64 ./cmd/testram
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./cmd/testram/testram-linux-arm64 ./cmd/testram

build:
	go build -o ./cmd/web/web ./cmd/web
	go build -o ./cmd/api/api ./cmd/api
	go build -o ./cmd/crawler/crawler ./cmd/crawler
	go build -o ./cmd/testcpu/testcpu ./cmd/testcpu
	go build -o ./cmd/testram/testram ./cmd/testram

builddocker:
	docker buildx build --platform=linux/amd64,linux/arm64 -t bootdotdev/synergychat-web ./cmd/web
	docker buildx build --platform=linux/amd64,linux/arm64 -t bootdotdev/synergychat-api ./cmd/api
	docker buildx build --platform=linux/amd64,linux/arm64 -t bootdotdev/synergychat-crawler ./cmd/crawler
	docker buildx build --platform=linux/amd64,linux/arm64 -t bootdotdev/synergychat-testram ./cmd/testram
	docker buildx build --platform=linux/amd64,linux/arm64 -t bootdotdev/synergychat-testcpu ./cmd/testcpu

pushdocker:
	docker push bootdotdev/synergychat-web
	docker push bootdotdev/synergychat-api
	docker push bootdotdev/synergychat-crawler
	docker push bootdotdev/synergychat-testram
	docker push bootdotdev/synergychat-testcpu
