FROM library/golang:1.10-alpine AS builder

RUN apk add --update --no-cache git ca-certificates make curl

WORKDIR /go/src/github.com/cmattoon/aws-ssm

COPY . .

RUN make deps

RUN make test-dkr

RUN make build-dkr

##
## Stage 2
##
FROM library/alpine:3.8

LABEL org.label-schema.schema-version = "1.0.0"
LABEL org.label-schema.version = "0.1.4"
LABEL org.label-schema.name = "aws-ssm"
LABEL org.label-schema.description = "Updates Kubernetes Secrets with AWS SSM Parameters"
LABEL org.label-schema.vendor = "com.cmattoon"
LABEL org.label-schema.vcs-url = "https://github.com/cmattoon/aws-ssm"

ENV AWS_REGION     ""
ENV AWS_ACCESS_KEY ""
ENV AWS_SECRET_KEY ""

ENV AWS_SSM_METRICS_URL "0.0.0.0:9999"
ENV AWS_SSM_INTERVAL    "60"

# Only required if running outside the cluster
ENV AWS_SSM_MASTER_URL  ""
ENV AWS_SSM_KUBE_CONFIG ""

COPY --from=builder /go/bin/aws-ssm /bin/aws-ssm
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER nobody

ENTRYPOINT ["/bin/aws-ssm"]
