# build
FROM golang:1-alpine as builder

RUN rm -rf /var/cache/apk/* && rm -rf /tmp/*
RUN apk update
RUN apk --no-cache add -U make git

WORKDIR /go/src/github.com/denouche/go-api-skeleton
COPY . /go/src/github.com/denouche/go-api-skeleton
RUN make deps build

# run
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/denouche/go-api-skeleton/go-api-skeleton .
CMD ["/go-api-skeleton"]

