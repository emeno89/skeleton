#example: docker build -f deploy/containers/cicd/prebuilder/Dockerfile --platform linux/amd64 -t prebuilder_amd64:1.20 .

FROM golangci/golangci-lint:latest-alpine

COPY proto/scripts/install.sh proto/scripts/install.sh

RUN apk update && \
    apk add --no-cache protobuf && \
    sh proto/scripts/install.sh
