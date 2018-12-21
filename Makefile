.PHONY:
AWS_REGION        ?= us-west-2
AWS_ACCESS_KEY    ?= none
AWS_SECRET_KEY    ?= none

# Docker Build
# ================================
DOCKERFILE_DIR     = .
DOCKERFILE         = Dockerfile

DOCKER_REPO       ?= cmattoon
IMAGE_NAME        ?= aws-ssm
IMAGE_TAG         ?= $(shell git describe --tags --always --dirty)

CURRENT_IMAGE      = $(DOCKER_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)
LATEST_IMAGE       = $(DOCKER_REPO)/$(IMAGE_NAME):latest

# Build/Output
# ================================
AWS_SSM_EXE        = build/aws-ssm
LDFLAGS           ?= -w -s

# Helm Chart
# ================================
RELEASE_NAME      ?= aws-ssm
RELEASE_NAMESPACE ?= kube-system
CHART_DIR         ?= $(IMAGE_NAME)
RBAC_ENABLED      ?= true
HOST_SSL_DIR      ?= ""
EXTRA_ARGS        ?= 


.PHONY: deps
deps:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	dep ensure -v

.PHONY: test
test:
	./scripts/go_test.sh

.PHONY: test-dkr
test-dkr:
	go test -v $(shell go list ./... | grep -v /vendor/)

.PHONY: build-dkr
build-dkr:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -v -a \
		-installsuffix ego \
		-ldflags="$(LDFLAGS)" \
		-o $(AWS_SSM_EXE)

.PHONY: build
build:
	CGO_ENABLED=0 \
		go build \
		-ldflags="$(LDFLAGS)" \
		-o $(AWS_SSM_EXE)

.PHONY: install
install:
	go install -v ./...

.PHONY: container
container:
	docker build -t $(CURRENT_IMAGE) $(DOCKERFILE_DIR) -f $(DOCKERFILE)
	docker tag $(CURRENT_IMAGE) $(LATEST_IMAGE)

.PHONY: chart
chart:
	helm lint aws-ssm

.PHONY: fmt
fmt:
	go fmt ./... -v

.PHONY: push-container
push-container: container
	docker push $(CURRENT_IMAGE)

.PHONY: install-chart
install-chart:
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
