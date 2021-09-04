.DEFAULT_GOAL := help

.PHONY:
AWS_REGION        ?= ap-southeast-2
AWS_ACCESS_KEY    ?=
AWS_SECRET_KEY    ?=

RELEASE_NAME      ?= aws-ssm
RELEASE_NAMESPACE ?= kube-system

DOCKER_REPO       ?= cmattoon
DOCKER_IMAGE      ?= aws-ssm
DOCKER_TAG        ?= $(strip $(shell git describe --tags --always --dirty))

GIT_REPO          ?= $(DOCKER_REPO)
GIT_PROJECT       ?= $(DOCKER_IMAGE)
GIT_URL           ?= https://github.com/$(GIT_REPO)/$(GIT_PROJECT)
COMMIT            ?= $(shell git log -1 --pretty=format:"%h")

CURRENT_IMAGE      = $(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG)
LATEST_IMAGE       = $(DOCKER_REPO)/$(DOCKER_IMAGE):latest

DOCKERFILE_DIR     = .
DOCKERFILE         = Dockerfile

# Output file
AWS_SSM_EXE        = build/aws-ssm-$(DOCKER_TAG)

CHART_DIR         ?= $(DOCKER_IMAGE)
RBAC_ENABLED      ?= true
HOST_SSL_DIR      ?=
ifeq ($(HOST_SSL_DIR),)
MOUNT_SSL=false
else
MOUNT_SSL=true
endif
EXTRA_ARGS        ?=

BUILD_DATE        ?= $(shell date +"%Y-%m-%dT%H:%M:%S")
BUILD_FLAGS       ?= -v
LDFLAGS           ?= -X github.com/cmattoon/aws-ssm/pkg/config.Version=$(DOCKER_TAG) -w -s

.PHONY: test
test: ## Runs Go tests
	./scripts/go_test.sh

.PHONY: dgoss
dgoss: ## Runs dgoss container tests
	dgoss run $(CURRENT_IMAGE)

build: ## Build the Go binary
build: $(AWS_SSM_EXE)
$(AWS_SSM_EXE):
	go build -o $(AWS_SSM_EXE) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)"

.PHONY: clean
clean: ## Clean files
clean:
	@rm build/*
	@find . -name '*~' -delete

.PHONY: container
container: ## Build the Docker image
container:
	docker build \
		--label org.label-schema.schema-version="1.0" \
		--label org.label-schema.name="$(DOCKER_REPO)/$(DOCKER_IMAGE)" \
		--label org.label-schema.description="Updates Kubernetes Secrets with AWS SSM Parameters" \
		--label org.label-schema.vendor="$(DOCKER_REPO)" \
		--label org.label-schema.build-date="$(BUILD_DATE)" \
		--label org.label-schema.vcs-url="$(GIT_URL)" \
		--label org.label-schema.vcs-ref="$(COMMIT)" \
		--label org.label-schema.version="$(COMMIT)" \
		-t $(CURRENT_IMAGE) $(DOCKERFILE_DIR) -f $(DOCKERFILE)
	docker tag $(CURRENT_IMAGE) $(LATEST_IMAGE)

.PHONY: chart
chart: ## Lint chart
chart:
	helm lint aws-ssm

.PHONY: push
push: ## Docker push
push:
	docker push $(CURRENT_IMAGE)

.PHONY: install
install: ## Install Helm Chart
install:
	helm upgrade --install $(RELEASE_NAME) \
		--namespace $(RELEASE_NAMESPACE) \
		--set image.tag=$(DOCKER_TAG) \
	 	--set aws.region=$(AWS_REGION) \
	 	--set aws.access_key=$(AWS_ACCESS_KEY) \
	 	--set aws.secret_key=$(AWS_SECRET_KEY) \
		--set rbac.enabled=$(RBAC_ENABLED) \
		--set ssl.mount_host=$(MOUNT_SSL) \
		--set ssl.host_path=$(HOST_SSL_DIR) \
	 	$(EXTRA_ARGS) $(CHART_DIR)

.PHONY: purge
purge: ## Purge Helm Chart
purge:
	helm del --purge $(RELEASE_NAME)

.PHONY: help
help: ## Show this message
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: login
login: ## Do a docker login
login:
	env | grep -i DOCKER | cut -d '=' -f1
	# docker login --username "${DOCKER_USER}" --password "${DOCKER_PASSWD}
