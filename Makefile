.PHONY:
AWS_REGION        ?= us-west-2
AWS_ACCESS_KEY    ?= none
AWS_SECRET_KEY    ?= none

RELEASE_NAME      ?= aws-ssm
RELEASE_NAMESPACE ?= kube-system

DOCKER_REPO       ?= cmattoon
IMAGE_NAME        ?= aws-ssm
IMAGE_TAG         ?= $(shell git log -1 --pretty=format:"%h")

CURRENT_IMAGE      = $(DOCKER_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)
LATEST_IMAGE       = $(DOCKER_REPO)/$(IMAGE_NAME):latest

DOCKERFILE_DIR     = .
DOCKERFILE         = Dockerfile

# Output file
AWS_SSM_EXE        = build/aws-ssm

CHART_DIR         ?= $(IMAGE_NAME)
RBAC_ENABLED      ?= true
HOST_SSL_DIR      ?= ""
EXTRA_ARGS        ?= 

.PHONY: test
test:
	./scripts/go_test.sh

.PHONY: build
build:
	go build -o $(AWS_SSM_EXE)

.PHONY: container
container:
	docker build -t $(CURRENT_IMAGE) $(DOCKERFILE_DIR) -f $(DOCKERFILE)
	docker tag $(CURRENT_IMAGE) $(LATEST_IMAGE)

.PHONY: chart
chart:
	helm lint aws-ssm

.PHONY: push-container
push-container: container
	docker push $(CURRENT_IMAGE)

.PHONY: install
install:
	helm upgrade --install $(RELEASE_NAME) \
		--namespace $(RELEASE_NAMESPACE) \
		--set image.tag=$(IMAGE_TAG) \
	 	--set aws.region=$(AWS_REGION) \
	 	--set aws.access_key=$(AWS_ACCESS_KEY) \
	 	--set aws.secret_key=$(AWS_SECRET_KEY) \
		--set rbac.enabled=$(RBAC_ENABLED) \
		--set host_ssl_dir=$(HOST_SSL_DIR) \
	 	$(EXTRA_ARGS) $(CHART_DIR)

.PHONY: purge
purge:
	helm del --purge $(RELEASE_NAME)
