FROM library/golang:1.10-alpine

RUN apk add --update --no-cache git

WORKDIR /go/src/github.com/cmattoon/aws-ssm

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

RUN go get -u -v github.com/kubernetes-sigs/aws-iam-authenticator/cmd/aws-iam-authenticator

## Stage 2
FROM library/alpine
LABEL org.label-schema.schema-version = "1.0.0"
LABEL org.label-schema.version = "0.1.5"
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

COPY --from=0 /go/bin/aws-iam-authenticator /bin/aws-iam-authenticator
COPY --from=0 /go/bin/aws-ssm /bin/aws-ssm

CMD ["aws-ssm"]
