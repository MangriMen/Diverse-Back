.PHONY: dev prod deploy testing

all:
	@echo "Usage: make BUILD_TARGET"
	@echo ""
	@echo "BUILD_TARGET:"
	@echo "\tdev\t\t-\tbuild with hot-reload"
	@echo "\tprod\t\t-\tproduction build"
	@echo "\tdeploy profile=\t-\tdeploy with profile prod or test"
	@echo "\testing\t\t-\trun base services for testing"

dev:
	docker-compose -p dev --profile dev up --build

prod:
	docker-compose -p prod --profile prod up --build

testing:
	sudo docker compose -p testing --profile testing stop postgres backend-testing

	sudo docker rm diverse-postgres-testing
	sudo docker rm diverse-backend-testing

	sudo docker volume rm testing_postgres-data

	sudo docker compose -p testing --profile testing up postgres backend-testing --build

deploy:
ifeq ($(profile),)
	(error profile argument not set)
endif

ifeq ($(BASE_HOST),)
	(BASE_HOST environment variable not set)
endif

ifeq ($(profile),prod)
	sed -i 's/:3040/:3030/g' docs/swagger.yml
endif

	sed -i 's/host: localhost/host: $(BASE_HOST)/g' docs/swagger.yml

	sudo docker compose -p $(profile) --profile $(profile) down
	sudo docker compose -p $(profile) --profile $(profile) pull
	sudo docker compose -p $(profile) --profile $(profile) up -d
