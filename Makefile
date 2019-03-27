SHELL := /bin/sh

MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(patsubst %/,%,$(dir $(MAKEFILE_PATH)))

DOCKER_IMAGE_NAME := $(if ${TRAVIS_REPO_SLUG},${TRAVIS_REPO_SLUG},supergiant/analyze)
NODEAGENT_DOCKER_IMAGE_NAME := $(if ${TRAVIS_REPO_SLUG},${TRAVIS_REPO_SLUG}-nodeagent,supergiant/analyze-nodeagent)
JOB_DOCKER_IMAGE_NAME := $(if ${TRAVIS_REPO_SLUG},${TRAVIS_REPO_SLUG}-registry-job,supergiant/analyze-registry-job)
DOCKER_IMAGE_TAG := $(if ${TAG},${TAG},$(shell git describe --tags --always | tr -d v || echo 'latest'))
GO111MODULE=on


define LINT
	@echo "Running code linters..."
	revive
	@echo "Running code linters finished."
endef

define GOIMPORTS
	goimports -v -w -local github.com/supergiant/analyze ${CURRENT_DIR}
endef

define TOOLS
		if [ ! -x "`which revive 2>/dev/null`" ]; \
        then \
        	echo "revive linter not found."; \
        	echo "Installing linter... into ${GOPATH}/bin"; \
        	GO111MODULE=off go get -u github.com/mgechev/revive ; \
        fi

        if [ ! -x "`which swagger 2>/dev/null`" ]; \
        then \
        	echo "swagger not found."; \
        	echo "Installing swagger... into ${GOPATH}/bin"; \
        	GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger ; \
        fi
endef


.PHONY: default
default: lint


.PHONY: lint
lint: tools
	@$(call LINT)

.PHONY: validate
validate: tools
	swagger validate ./swagger/api-spec.yml

.PHONY: gen-swagger
gen-swagger: validate
	swagger generate model \
    		--target=./pkg \
    		--spec=./swagger/api-spec.yml
	swagger generate server \
		--target=./pkg \
		--server-package=api \
		--spec=./swagger/api-spec.yml \
		--exclude-main \
		--name=analyze \
		--existing-models=github.com/supergiant/analyze/pkg/models
	cp ./swagger/api-spec.yml ./asset/swagger/api-spec.yml

.PHONY: test
test:
	go test -mod=vendor -count=1 -tags=dev -race ./...

.PHONY: tools
tools:
	@$(call TOOLS)

.PHONY: goimports
goimports:
	@$(call GOIMPORTS)

.PHONY: build-image
build-image:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .
	docker tag $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) $(DOCKER_IMAGE_NAME):latest

	docker build -t $(NODEAGENT_DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) -f cmd/analyze-nodeagent/Dockerfile .
	docker tag $(NODEAGENT_DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) $(NODEAGENT_DOCKER_IMAGE_NAME):latest

	docker build -t $(JOB_DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) -f cmd/analyze-registry-job/Dockerfile .
	docker tag $(JOB_DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) $(JOB_DOCKER_IMAGE_NAME):latest

.PHONY: push
push:
	docker push $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
	docker push $(NODEAGENT_DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
	docker push $(JOB_DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

.PHONY: gofmt
gofmt:
	go fmt ./...

.PHONY: fmt
fmt: gofmt goimports


#all dependencies are "go get"-able for general dev environment usability.
# To compile all protobuf files in this repository, run "make protobuf"
.PHONY: gen-protobuf
gen-protobuf:
	docker run --rm -v ${CURRENT_DIR}/pkg/plugin:/defs namely/protoc-all:latest -i proto -l go -d /defs -o .

.PHONY: gen-assets
gen-assets:
	docker run --rm -it --name supergiant_frontend_builder \
		--mount type=bind,src=${CURRENT_DIR},dst=/tmp \
		-w /usr/src/app node:10-alpine \
		sh -c "cp -a /tmp/ui/. /usr/src/app && ls -la && npm i && npm run build && cp -a /usr/src/app/dist/. /tmp/asset/ui"
	cd ${CURRENT_DIR}/asset && go generate -mod=vendor
