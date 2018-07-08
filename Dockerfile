FROM library/golang:alpine

RUN apk add --no-cache git

WORKDIR /go/src/github.com/cmattoon/aws-ssm

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["aws-ssm"]
