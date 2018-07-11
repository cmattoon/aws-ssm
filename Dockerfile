FROM library/golang:alpine

LABEL org.label-schema.schema-version = "1.0.0"
LABEL org.label-schema.name = "aws-ssm"
LABEL org.label-schema.description = "Updates Kubernetes Secrets with AWS SSM Parameters"
LABEL org.label-schema.vendor = "com.cmattoon"
LABEL org.label-schema.vcs-url = "https://github.com/cmattoon/aws-ssm"

RUN apk add --no-cache git

WORKDIR /go/src/github.com/cmattoon/aws-ssm

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["aws-ssm"]

