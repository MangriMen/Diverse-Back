.PHONY: dev prod

dev:
	docker-compose -f dev.docker-compose.yml up --build

prod:
	docker-compose -f build.docker-compose.yml up --build