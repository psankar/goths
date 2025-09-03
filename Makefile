.PHONY: all generate-templ generate-sqlc
all: generate-templ generate-sqlc
	@if [ -f .templ_success ] && [ -f .sqlc_success ]; then \
		echo "Launching Goths server ..."; \
		rm -f .templ_success .sqlc_success; \
		go run main.go; \
	else \
		echo "One or both commands failed, skipping go run main.go"; \
		rm -f .templ_success .sqlc_success; \
		exit 1; \
	fi

generate-templ:
	@cd templ && templ generate && touch ../.templ_success || true

generate-sqlc:
	@cd sqlc && sqlc generate && touch ../.sqlc_success || true