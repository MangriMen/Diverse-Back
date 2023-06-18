.PHONY: dev test deploy

COMPOSE=docker compose

define gen_compose_flags
	--profile $(1)
endef

define gen_compose_command
	PROFILE=$(strip $(1)) $(COMPOSE) $(call gen_compose_flags,$(1))
endef

COMPOSE_DEPLOY_COMMAND=$(call gen_compose_command,$(profile))

UP_FLAGS_DEV=--build
UP_FLAGS_DEPLOY=-d --force-recreate

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
	@echo "\ttest\t\t-\tproduction test build"
	@echo "\tdev\t\t-\tbuild with hot-reload"
	@echo "\tdeploy profile=\t-\tdeploy with profile prod or test"

dev:
	$(call gen_compose_command,$(DEV_PROFILE)) up $(UP_FLAGS_DEV)

test:
	$(call gen_compose_command,$(TEST_PROFILE)) up $(UP_FLAGS_DEV)

deploy: prepare_swagger_to_deploy
	$(COMPOSE_DEPLOY_COMMAND) pull
	$(COMPOSE_DEPLOY_COMMAND) up $(UP_FLAGS_DEPLOY)

prepare_swagger_to_deploy: check_deploy_environment
ifeq ($(profile),$(TEST_PROFILE))
	sed -i 's/$(DEV_HOST)/$(PROD_HOST)/g' $(SWAGGER_DOC)
endif

ifeq ($(profile),$(DEV_PROFILE))
	sed -i 's/$(DEV_HOST)/$(PROD_HOST)/g' $(SWAGGER_DOC)
endif

check_deploy_environment: check_profile_exists check_base_host_exists

check_profile_exists:
ifeq ($(profile),)
	(error profile argument not set)
endif

check_base_host_exists:
ifeq ($(profile),$(TEST_PROFILE))
ifeq ($(BASE_HOST),)
	(BASE_HOST environment variable not set)
endif
endif

