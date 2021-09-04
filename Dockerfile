###
## Stage I - Build aws-ssm binary
#
FROM library/golang:1.14-alpine

RUN apk add --update --no-cache git

WORKDIR /go/src/github.com/cmattoon/aws-ssm

COPY . .

RUN go install -v ./...

###
## Stage II - Install aws-iam-authenticator
#
FROM library/golang:1.14-alpine

RUN apk add --update --no-cache git

RUN go get -u -v sigs.k8s.io/aws-iam-authenticator/cmd/aws-iam-authenticator


###
## Stage III - Add ca-certificates, binaries
#
FROM library/alpine:3.11

ENV AWS_REGION     ""
ENV AWS_ACCESS_KEY ""
ENV AWS_SECRET_KEY ""
ENV METRICS_URL    "0.0.0.0:9999"

# Only required if running outside the cluster
ENV MASTER_URL     ""
ENV KUBE_CONFIG    ""

RUN apk add --update ca-certificates


COPY --from=1 /go/bin/aws-iam-authenticator /bin/aws-iam-authenticator
COPY --from=0 /go/bin/aws-ssm /bin/aws-ssm

ENTRYPOINT ["/bin/aws-ssm"]
