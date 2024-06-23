build:
	@go build -o bin/gogym

run: build
	@./bin/gogym

test:
	@go test -v ./...