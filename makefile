build:
	@go build -o bin/api ./cmd/api
run: build
	@./bin/api
continue:
	@air

# NOTE: when use the make file to run, always build and then run.