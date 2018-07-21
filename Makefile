.PHONY:
DOCKER_REPO=cmattoon
IMAGE_NAME=aws-ssm
IMAGE_TAG:=$(shell git log -1 --pretty=format:"%h")

CURRENT_IMAGE=$(DOCKER_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)
LATEST_IMAGE=$(DOCKER_REPO)/$(IMAGE_NAME):latest

DOCKERFILE_DIR=.
DOCKERFILE=Dockerfile

AWS_SSM_EXE=build/aws-ssm

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
	docker push $(LATEST_IMAGE)
