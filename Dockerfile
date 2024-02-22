###
## Stage I - Build aws-ssm binary
#
FROM library/golang:1.16-alpine

RUN apk add --update --no-cache git

WORKDIR /go/src/github.com/cmattoon/aws-ssm

COPY . .

RUN go install -v ./...

###
## Stage II - Install aws-iam-authenticator
#
FROM library/alpine:3.14

WORKDIR /tmp
RUN wget https://amazon-eks.s3.us-west-2.amazonaws.com/1.21.2/2021-07-05/bin/linux/amd64/aws-iam-authenticator
RUN chmod +x aws-iam-authenticator


###
## Stage III - Add ca-certificates, binaries
#
FROM library/alpine:3.14

ENV AWS_REGION     ""
ENV AWS_ACCESS_KEY ""
ENV AWS_SECRET_KEY ""
ENV METRICS_URL    "0.0.0.0:9999"

# Only required if running outside the cluster
ENV MASTER_URL     ""
ENV KUBE_CONFIG    ""

RUN apk add --update ca-certificates


COPY --from=1 /tmp/aws-iam-authenticator /bin/aws-iam-authenticator
COPY --from=0 /go/bin/aws-ssm /bin/aws-ssm

ENTRYPOINT ["/bin/aws-ssm"]
