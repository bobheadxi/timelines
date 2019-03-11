
all:
	go build

.PHONY: web
web:
	cd web && npm run start

.PHONY: devenv
devenv:
	docker-compose -f dev/docker-compose.yml up -d
	docker ps

.PHONY: devmonitoring
devmonitoring:
	docker-compose -f dev/monitoring.yml up -d
	docker ps

.PHONY: clean
clean:
	docker-compose -f dev/monitoring.yml rm -f -s -v
	docker-compose -f dev/docker-compose.yml rm -f -s -v
	rm -rf tmp

.PHONY: graphql
graphql:
	$(MAKE) -C graphql go
	$(MAKE) -C graphql ts

.PHONY: devpg
devpg: pg-reset pg-init

.PHONY: pg-reset
pg-reset:
	docker exec -i postgres psql -U bobheadxi projector-dev < db/sql/reset.sql

.PHONY: pg-init
pg-init:
	docker exec -i postgres psql -U bobheadxi projector-dev < db/sql/repos.sql
