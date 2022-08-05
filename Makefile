deps:
	go get ./...

test:
	go test -v ./...

build:
	docker build -t rami-notion-api --build-arg NOTION_API_KEY .

deploy:
	fly deploy \
        --build-secret NOTION_API_KEY=${NOTION_API_KEY}

run:
	go run ./cmd/api/main.go

fmt:
	go fmt ./...

check:
	go vet ./...

dev:
	reflex -c reflex.conf

clean:
	@echo "clearing test cache"
	go clean -testcache

.PHONY: deps test build deploy run fmt check dev clean
