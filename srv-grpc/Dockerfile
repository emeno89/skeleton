FROM golang:1.20-alpine AS builder

ARG NAME=srv-grpc
ARG MODE=api

RUN apk update && apk add --no-cache ca-certificates protobuf

WORKDIR ${GOPATH}/src/skeleton

COPY go.sum go.mod ./
RUN go mod download

COPY proto/scripts proto/scripts

RUN sh proto/scripts/install.sh

COPY shared shared
COPY proto/schema proto/schema

RUN sh proto/scripts/generate.sh

WORKDIR ${GOPATH}/src/skeleton/${NAME}

COPY ${NAME}/internal/${MODE} ./internal/${MODE}
COPY ${NAME}/internal/pkg ./internal/pkg
COPY ${NAME}/cmd/${MODE} ./cmd/${MODE}

WORKDIR ${GOPATH}/src/skeleton/${NAME}/cmd/${MODE}

RUN GOOS=linux GOARCH=amd64 go build -o /go/bin/svc

FROM alpine:latest

ARG NAME=srv-grpc

RUN apk update && apk add --no-cache curl tzdata

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/svc /go/bin/svc

#i18n files
COPY ${NAME}/cmd/translation /go/bin/translation

WORKDIR /go/bin

ENTRYPOINT ["/go/bin/svc"]
