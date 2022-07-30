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

fmt:
	go fmt ./...

check:
	go vet ./...

.PHONY: deps test build deploy run fmt check
