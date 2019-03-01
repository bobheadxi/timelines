
.PHONY: web
web:
	cd web && npm run start

.PHONY: devenv
devenv:
	docker-compose -f dev/docker-compose.yml up -d
	docker ps

.PHONY: clean
clean:
	docker-compose -f dev/docker-compose.yml down
