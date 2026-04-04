.PHONY: run build migrate test clean

# Run dev
run:
	go run cmd/server/main.go

# Build
build:
	go build -o bin/aulaflash cmd/server/main.go

# Run migrations
migrate:
	go run cmd/migrate/main.go

# Test
test:
	go test -v ./...

# Clean
clean:
	rm -rf bin/
	rm -rf /tmp/aulaflash-uploads
