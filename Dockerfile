FROM golang:1.7.1-alpine

RUN apk add --update git && apk add --update make && rm -rf /var/cache/apk/*

WORKDIR /go/src/github.com/${GITHUB_ORG:-ernestio}/api-gateway

COPY Makefile /go/src/github.com/${GITHUB_ORG:-ernestio}/api-gateway
RUN make deps

COPY . /go/src/github.com/${GITHUB_ORG:-ernestio}/api-gateway
RUN go install

ENTRYPOINT ./entrypoint.sh
