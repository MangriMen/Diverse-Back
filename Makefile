.PHONY: dev test deploy undeploy

COMPOSE=docker compose

UP_FLAGS_DEV=--build
UP_FLAGS_DEPLOY=-d

DEV_HOST=host: localhost
PROD_HOST=host: $(BASE_HOST)

DEV_PORT=3040
PROD_PORT=3030

DEV_PROFILE=dev
TEST_PROFILE=test
PROD_PROFILE=prod

SWAGGER_DOC=docs/swagger.yml

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

deploy: undeploy prepare_swagger_to_deploy
	# PROFILE=$(profile) $(COMPOSE) --profile $(profile) pull
	# PROFILE=$(profile) $(COMPOSE) --profile $(profile) up $(UP_FLAGS_DEPLOY)

undeploy: check_deploy_environment
	PROFILE=$(profile) $(COMPOSE) --profile $(profile) down

prepare_swagger_to_deploy: check_deploy_environment
ifeq ($(profile),$(PROD_PROFILE))
	sed -i 's/:$(DEV_PORT)/:$(PROD_PORT)/g' $(SWAGGER_DOC)
endif

	sed -i 's/$(DEV_HOST)/$(PROD_HOST)/g' $(SWAGGER_DOC)

check_deploy_environment:
ifeq ($(profile),)
	(error profile argument not set)
endif

ifeq ($(BASE_HOST),)
	(BASE_HOST environment variable not set)
endif