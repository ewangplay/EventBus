#Copyright ArxanFintech Technology Ltd. 2017-2018 All Rights Reserved.
#
#Licensed under the Apache License, Version 2.0 (the "License");
#you may not use this file except in compliance with the License.
#You may obtain a copy of the License at
#
#		 http://www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS,
#WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#See the License for the specific language governing permissions and
#limitations under the License.
#
# -------------------------------------------------------------
# This makefile defines the following targets
#
#   - all (default) - builds all targets and runs all tests/checks
#   - checks - runs all tests/checks
#   - srvc - builds services binary
#   - unit-test - runs the go-test based unit tests
#   - behave - runs the behave test
#   - behave-deps - ensures pre-requisites are availble for running behave manually
#   - gotools - installs go tools like golint
#   - linter - runs all code checks
#   - clean - cleans the build area
#   - dist-clean - superset of 'clean' that also removes persistent state

## Project Parameters
ORG_NAME=ewangplay
DOCKER_NS=ewangplay
PROJECT_NAME=$(ORG_NAME)/eventbus
PKG_NAME = github.com/$(PROJECT_NAME)
BASE_VERSION = 1.0.0
IS_RELEASE = false

## Docker Images Parameters
EVENTBUS_IMAGE_NAME = eventbus

## Get system info
UID = $(shell id -u)
ARCH=$(shell uname -m)

## Build Project Version
ifneq ($(IS_RELEASE),true)
EXTRA_VERSION ?= snapshot-$(shell git rev-parse --short HEAD)
PROJECT_VERSION=$(ARCH)-$(BASE_VERSION)-$(EXTRA_VERSION)
else
PROJECT_VERSION=$(ARCH)-$(BASE_VERSION)
endif

EXECUTABLES = go git
K := $(foreach exec,$(EXECUTABLES),\
	$(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH: Check dependencies")))

# SUBDIRS are components that have their own Makefiles that we can invoke
SUBDIRS = gotools
SUBDIRS:=$(strip $(SUBDIRS))

PROJECT_FILES = $(shell git ls-files)

all: srvc checks

checks: linter unit-test #behave

srvc: eventbus

.PHONY: $(SUBDIRS)
$(SUBDIRS):
	cd $@ && $(MAKE)

.PHONY: eventbus
eventbus: build/bin/eventbus

.PHONY: docker
docker-deps: eventbus
docker: docker-deps
	docker build -t $(DOCKER_NS)/$(EVENTBUS_IMAGE_NAME) -f docker/eventbus/Dockerfile ./
	docker tag $(DOCKER_NS)/$(EVENTBUS_IMAGE_NAME) $(DOCKER_NS)/$(EVENTBUS_IMAGE_NAME):$(PROJECT_VERSION)

.PHONY: docker-clean
docker-clean:
	docker images -q $(DOCKER_NS)/$(EVENTBUS_IMAGE_NAME) | uniq | xargs -I '{}' docker rmi -f '{}'

unit-test: gotools
	@./scripts/goUnitTests.sh

bench: gotools
	@./scripts/goBenchTests.sh

behave-deps: eventbus
behave: behave-deps
	@echo "Running behave tests"
	@cd bddtests; behave $(BEHAVE_OPTS)

linter: gotools
	@echo "LINT: Running code checks.."
	@echo "Running go vet"
	go vet ./adapter/...
	go vet ./common/...
	go vet ./config/...
	go vet ./driver/...
	go vet ./i/...
	go vet ./log/...
	go vet ./rest/...
	go vet ./services/...
	go vet ./utils/...
	@echo "Running goimports"
	@./scripts/goimports.sh

build/bin:
	mkdir -p $@

build/bin/%:
	@mkdir -p $(@D)
	@echo "$@"
	GOBIN=$(abspath $(@D)) go install $(PKG_NAME)/services/$(@F)
	@echo "Binary available as $@"
	@touch $@

.PHONY: $(SUBDIRS:=-clean)
$(SUBDIRS:=-clean):
	cd $(patsubst %-clean,%,$@) && $(MAKE) clean

.PHONY: clean
clean: docker-clean
	-@rm -rf build ||:

.PHONY: dist-clean
dist-clean: clean gotools-clean
	-@rm -rf /opt/$(ORG_NAME)/eventbus/* ||:
