deps:
	go get ./...

test:
	go test -v ./...

build:
	docker build -t rami-notion-api .

deploy:
	echo "Deploying..."

run:
	go run ./cmd/api/main.go

.PHONY: deps test build deploy run