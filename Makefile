
build:
	@go build

container:
	docker build -t cmattoon/aws-ssm .
