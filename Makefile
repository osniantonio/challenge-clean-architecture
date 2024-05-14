build:
	docker-compose up -d --build
	@echo "Contêineres do aplicativo compilados e iniciados com sucesso."

up:
	docker-compose up -d
	@echo "Contêineres do aplicativo iniciados com sucesso."

down:
	docker-compose down
	@echo "Contêineres do aplicativo parados e removidos com sucesso."

logs:
	docker-compose logs mysql

.PHONY: go