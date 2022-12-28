lint_docker_compose_file = "./development/golangci_lint/docker-compose.yml"

lint-build:
	@echo "🌀 ️container are building..."
	@docker-compose --file=$(lint_docker_compose_file) build -q
	@echo "✔  ️container built"

lint-check:
	@echo "🌀️ code linting..."
	@docker-compose --file=$(lint_docker_compose_file) run gin-golinter golangci-lint run \
 		&& echo "✔️  checked without errors" \
 		|| echo "☢️  code style issues found"


lint-fix:
	@echo "🌀 ️code fixing..."
	@docker-compose --file=$(lint_docker_compose_file) run gin-golinter golangci-lint run --fix \
		&& echo "✔️  fixed without errors" \
		|| (echo "⚠️️  you need to fix above issues manually" && exit 1)
	@echo "⚠️️ run \"make lint-check\" again to check what did not fix yet"
db-migrate:
	goose -dir "./server/db/migrations" ${DB_DRIVER} "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" up

swag-gen:
	swag init