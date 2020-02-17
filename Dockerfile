###
## Stage I - Build aws-ssm binary
#
FROM library/golang:1.13-alpine

RUN apk add --update --no-cache git

WORKDIR /go/src/github.com/cmattoon/aws-ssm

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

###
## Stage II - Add ca-certificates, binaries, aws-iam-authenticator
#
FROM library/alpine:3.7

ENV AWS_REGION     ""
ENV AWS_ACCESS_KEY ""
ENV AWS_SECRET_KEY ""
ENV METRICS_URL    "0.0.0.0:9999"

# Only required if running outside the cluster
ENV MASTER_URL     ""
ENV KUBE_CONFIG    ""

RUN apk add --update ca-certificates

ADD https://amazon-eks.s3-us-west-2.amazonaws.com/1.14.6/2019-08-22/bin/linux/amd64/aws-iam-authenticator /bin
COPY --from=0 /go/bin/aws-ssm /bin/aws-ssm

ENTRYPOINT ["/bin/aws-ssm"]
