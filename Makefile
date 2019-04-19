COMMIT=`git rev-parse HEAD`

all:
	@echo Version $(COMMIT)
	go build -ldflags "-X github.com/bobheadxi/timelines/config.Commit=$(COMMIT)"

.PHONY: clean
clean:
	docker-compose -f dev/monitoring.yml rm -f -s -v
	docker-compose -f dev/docker-compose.yml rm -f -s -v
	rm -rf tmp

.PHONY: web
web:
	cd web && npm run build

.PHONY: scripts
scripts:
	$(MAKE) -C .scripts install

.PHONY: lint
lint:
	./.scripts/lint.sh

# Dev environment

.PHONY: devenv
devenv:
	docker-compose -f dev/docker-compose.yml up -d
	docker ps

.PHONY: devmonitoring
devmonitoring:
	docker-compose -f dev/monitoring.yml up -d
	docker ps

# Codegen

.PHONY: graphql
graphql:
	$(MAKE) -C graphql go

# PG utils

.PHONY: devpg
devpg: pg-reset pg-init

.PHONY: pg-reset
pg-reset:
	docker exec -i postgres psql -U bobheadxi timelines-dev < db/sql/reset.sql

.PHONY: pg-init
pg-init:
	docker exec -i postgres psql -U bobheadxi timelines-dev < db/sql/repos.sql

.PHONY: herokupg
herokupg:
	heroku pg:psql --app timelines-server < db/sql/reset.sql
	heroku pg:psql --app timelines-server < db/sql/repos.sql

GOOGLE_APPLICATION_CREDENTIALS_RAW=`< gcp.json`
.PHONY: herokugcp
herokugcp:
	heroku config:set GOOGLE_APPLICATION_CREDENTIALS_RAW="$(GOOGLE_APPLICATION_CREDENTIALS_RAW)"
