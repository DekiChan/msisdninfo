FROM golang:alpine

COPY . /go/src/github.com/dekichan/msisdninfo
WORKDIR /go/src/github.com/dekichan/msisdninfo

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -o msisdninfo .
CMD ["msisdninfo"]
EXPOSE 8080
