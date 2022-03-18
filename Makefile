.PHONY: help build.web start stop clean shell docs lint

default: help

build.web: ## Build the web container
	@docker-compose build course-manager

help: ## Show this help
	@echo
	@fgrep -h " ## " $(MAKEFILE_LIST) | fgrep -v fgrep | sed -Ee 's/([a-z.]*):[^#]*##(.*)/\1##\2/' | column -t -s "##"
	@echo

start: ## Run the application locally in the background
	@docker compose up --build --detach course-manager

stop: ## Stop the application
	@docker compose stop

clean: ## Remove containers and delete all data from the local volumes
	@docker compose down --remove-orphans --volumes

shell: ## Shell into a development container
	@docker compose exec course-manager sh

logs: ## Show the application logs
	@docker compose logs --follow --tail=1000 course-manager

run: start logs ## Run the application locally

docs: ## Run the documentation service
	@docker compose up --build --detach swagger-ui
	@echo "Swagger running: http://localhost:8080"
	@echo

lint: ## Lint the application code
	@docker compose run --rm lint
	@echo "Lint completed with no errors."