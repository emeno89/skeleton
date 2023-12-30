#!/bin/sh

(
  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && \
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
) || false