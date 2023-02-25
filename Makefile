.PHONY: dev prod

all:
	@echo "Usage: make BUILD_TARGET"
	@echo ""
	@echo "BUILD_TARGET:"
	@echo "\tdev\t-\tbuild with hot-reload"
	@echo "\tprod\t-\tproduction build"

dev:
	docker-compose --profile dev up --build

prod:
	docker-compose --profile prod up --build