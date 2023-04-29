.PHONY: dev test deploy testing

COMPOSE=docker compose

UP_FLAGS_DEV=--build
UP_FLAGS_DEPLOY=-d

DEV_PROFILE=dev
TEST_PROFILE=test
PRODUCTION_PROFILE=prod

all:
	@echo "Usage: make BUILD_TARGET"
	@echo ""
	@echo "BUILD_TARGET:"
	@echo "\tdev\t\t-\tbuild with hot-reload"
	@echo "\tprod\t\t-\tproduction build"
	@echo "\tdeploy profile=\t-\tdeploy with profile prod or test"
	@echo "\testing\t\t-\trun base services for testing"

dev:
	PROFILE=$(DEV_PROFILE) $(COMPOSE) --profile $(DEV_PROFILE) up $(UP_FLAGS_DEV)

test:
	PROFILE=$(TEST_PROFILE) $(COMPOSE) --profile $(TEST_PROFILE) up $(UP_FLAGS_DEV)

deploy:
ifeq ($(profile),)
	(error profile argument not set)
endif

ifeq ($(BASE_HOST),)
	(BASE_HOST environment variable not set)
endif

ifeq ($(profile),$(PRODUCTION_PROFILE))
	sed -i 's/:3040/:3030/g' docs/swagger.yml
endif

	sed -i 's/host: localhost/host: $(BASE_HOST)/g' docs/swagger.yml

	PROFILE=$(profile) $(COMPOSE) --profile $(profile) down
	PROFILE=$(profile) $(COMPOSE) --profile $(profile) pull
	PROFILE=$(profile) $(COMPOSE) --profile $(profile) up $(UP_FLAGS_DEPLOY)
