.PHONY: dev prod deploy-prod deploy-test

all:
	@echo "Usage: make BUILD_TARGET"
	@echo ""
	@echo "BUILD_TARGET:"
	@echo "\tdev\t-\tbuild with hot-reload"
	@echo "\tprod\t-\tproduction build"
	@echo "\deploy profile=\t-\tdeploy with profile prod or test"

dev:
	docker-compose --profile dev up --build

prod:
	docker-compose --profile prod up --build

BASE_DIR=diverse

deploy:
	cd ~/$(BASE_DIR)

	git clone https://$TOKEN@github.com/$REPO.git $REPO-$(profile)
	cd $REPO-$(profile)
	git checkout develop
	git pull

	@if [ $(profile) = "prod" ]; then\
        sed -i 's/:3040/:3030/g' docs/swagger.yml
    fi

	sudo docker compose -p $(profile) --profile $(profile) down
	sudo docker compose -p $(profile) --profile $(profile) pull
	sudo docker compose -p $(profile) --profile $(profile) up -d
