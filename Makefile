SHELL = /usr/bin/bash

APP_NAME := go-news-comments

$(eval TAGVERSION := $(shell git describe --tags))
$(eval HASHCOMMIT := $(shell git log --pretty=tformat:"%h" -n1 ))
$(eval BRANCHNAME := $(shell git branch --show-current))
ifeq ($(TAGVERSION),undefined)
    # default tag is undefined
    VERSION := $(BRANCHNAME)
else ifeq ($(TAGVERSION),)
    # is empty tag 
    VERSION := $(BRANCHNAME)
else
    VERSION := $(TAGVERSION)
endif
$(eval VERSIONDATE := $(shell git show -s --format=%cI $($VERSION)))

PG_STARTED=$(shell echo $$((`docker ps --filter "name=db-comments" --quiet 2> /dev/null | wc -l` + `ps aux|grep -m 1 [p]ostgres:| wc -l`+0)))
pg-run:
ifeq ($(PG_STARTED),0)
	docker run --name db-comments -d -e POSTGRES_HOST_AUTH_METHOD=trust -p 5433:5432 postgres:15.4
	sleep 2
	psql -d 'postgres://postgres@localhost:5433/postgres?sslmode=disable' -c "CREATE DATABASE comments;"
	psql -d 'postgres://postgres@localhost:5433/comments?sslmode=disable' -f scripts/schema.sql
endif

PG_TEST_STARTED=$(shell echo $$((`docker ps --filter "name=db-comments-test" --quiet 2> /dev/null | wc -l` + `ps aux|grep -m 1 [p]ostgres:| wc -l`+0)))
pg-run-test:
ifeq ($(PG_TEST_STARTED),0)
	docker run --name db-comments-test -d -e POSTGRES_HOST_AUTH_METHOD=trust -p 5433:5432 postgres:15.4
	sleep 2
	psql -d 'postgres://postgres@localhost:5433/postgres?sslmode=disable' -c "CREATE DATABASE comments;"
	psql -d 'postgres://postgres@localhost:5433/comments?sslmode=disable' -f scripts/schema.sql
	#psql -d 'postgres://postgres@localhost:5432/comments?sslmode=disable' -f scripts/test_fixtures.sql
endif

build:
	@go mod tidy && go build -ldflags="-X 'github.com/mstyushin/go-news-comments/pkg/config.Version=$(VERSION)' -X 'github.com/mstyushin/go-news-comments/pkg/config.Hash=$(HASHCOMMIT)' -X 'github.com/mstyushin/go-news-comments/pkg/config.VersionDate=$(VERSIONDATE)'" -o bin/$(APP_NAME) github.com/mstyushin/go-news-comments/cmd/server
	@chmod +x bin/$(APP_NAME)

run: build pg-run
	@mkdir -p bin log
	@bin/$(APP_NAME) > log/$(APP_NAME).log 2>&1 & echo "$$!" > /tmp/$(APP_NAME).pid

test: pg-run-test
	@go mod tidy && go test -v ./...

stop:
	-pkill -f $(APP_NAME)

clean: stop
	@rm -f bin/*
	@rm -f log/*
	@rm -f /tmp/$(APP_NAME).pid
	docker stop db-comments || docker stop db-comments-test || true
	docker container prune -f && docker volume prune -f
