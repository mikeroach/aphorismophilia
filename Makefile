all: lint test

lint:
	 goreportcard-cli -v -t 100 .
	 golangci-lint run -v --enable-all ./... 

deeplint: 
	golangci-lint run -v --exclude-use-default=false --enable-all ./...

test:
	CGO_ENABLED=0 go test -v ./...

test-coverage:
	go test -coverprofile=/tmp/a13a-testcoverage.out ./...

build:
	go build -v ./...

run:
	go run -v .

docker-test:
	docker build --target test -t aphorismophilia-test:local .

docker-release:
	docker build --target release -t aphorismophilia:local .

docker-run:
	docker run -p 8888:8888 --rm --name aphorismophilia aphorismophilia:local
