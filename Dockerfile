FROM library/golang:1.10-alpine

RUN apk add --update --no-cache git

WORKDIR /go/src/github.com/cmattoon/aws-ssm

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/aws-ssm .

## Stage 2
FROM library/alpine

ARG TAG_VERSION

LABEL org.label-schema.schema-version = "1.0.0"
LABEL org.label-schema.version = "$TAG_VERSION"
LABEL org.label-schema.name = "aws-ssm"
LABEL org.label-schema.description = "Updates Kubernetes Secrets with AWS SSM Parameters"
LABEL org.label-schema.vendor = "com.cmattoon"
LABEL org.label-schema.vcs-url = "https://github.com/cmattoon/aws-ssm"

ENV AWS_REGION     ""
ENV AWS_ACCESS_KEY ""
ENV AWS_SECRET_KEY ""
ENV METRICS_URL    "0.0.0.0:9999"

# Only required if running outside the cluster
ENV MASTER_URL     ""
ENV KUBE_CONFIG    ""

RUN apk add --update ca-certificates

COPY --from=0 /go/bin/aws-ssm /bin/aws-ssm

CMD ["aws-ssm"]
