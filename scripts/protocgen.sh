#!/usr/bin/env bash

set -eo pipefail

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
    protoc \
    -I "proto" \
    -I "third_party/proto" \
    --gocosmos_out=plugins=interfacetype+grpc:. \
    $(find "${dir}" -maxdepth 1 -name '*.proto')
    
    # command to generate gRPC gateway (*.pb.gw.go in respective modules) files
    protoc \
    -I "proto" \
    -I "third_party/proto" \
    --grpc-gateway_out=logtostderr=true:. \
    $(find "${dir}" -maxdepth 1 -name '*.proto')
    
done

# move proto files to the right places
cp -r github.com/cosmos/ethermint/* ./
rm -rf github.com
