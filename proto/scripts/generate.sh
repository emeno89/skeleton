#!/bin/sh

rm -rf shared/grpc/pb || true

cd proto/schema || false

for f in *.proto
do
  protoc "${f}" --go_out=../../../ --go-grpc_out=../../../;
done