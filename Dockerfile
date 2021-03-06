FROM golang:1.10-alpine as builder
RUN apk add --update make git
WORKDIR /go/src/git.containerum.net/ch/resource-service
COPY . .
RUN VERSION=$(git describe --abbrev=0 --tags) make build-for-docker

FROM alpine:3.7

COPY --from=builder /tmp/resource /
ENV CH_RESOURCE_DEBUG="true" \
    CH_RESOURCE_TEXTLOG="true" \
    CH_RESOURCE_MONGO_ADDR="http://mongo:27017" \
    CH_RESOURCE_MIN_SERVICE_PORT=30000 \
    CH_RESOURCE_MAX_SERVICE_PORT=32767
EXPOSE 1213

CMD "/resource"
