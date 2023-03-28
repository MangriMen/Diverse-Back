.PHONY: dev prod deploy

all:
	@echo "Usage: make BUILD_TARGET"
	@echo ""
	@echo "BUILD_TARGET:"
	@echo "\tdev\t\t-\tbuild with hot-reload"
	@echo "\tprod\t\t-\tproduction build"
	@echo "\tdeploy profile=\t-\tdeploy with profile prod or test"

dev:
	docker-compose --profile dev up --build

prod:
	docker-compose --profile prod up --build

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
