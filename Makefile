.PHONY: all generate-templ generate-sqlc docker
all: docker
	@echo "Launching Goths server using docker-compose..."
	docker-compose up

docker:
	@echo "Generating SQL files..."
	@docker run --rm -v $(PWD)/sqlc:/src -w /src sqlc/sqlc:latest generate
	@echo "Generating template files..."
	@docker run --rm -v $(PWD):/src -w /src/templ ghcr.io/a-h/templ:latest generate
	@echo "Building Docker image..."
	@docker build -t goths:latest .
