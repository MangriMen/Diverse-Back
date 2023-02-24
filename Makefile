.PHONY: dev prod

all:
	@echo "Usage: make BUILD_TARGET"
	@echo ""
	@echo "BUILD_TARGET:"
	@echo "\tdev\t-\tbuild with hot-reload"
	@echo "\tprod\t-\tproduction build"

dev:
	docker-compose -f dev.docker-compose.yml up --build

prod:
	docker-compose -f prod.docker-compose.yml up --build