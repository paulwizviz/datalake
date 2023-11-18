#!/bin/bash

## Block specification
protoc --proto_path=$PWD/protos \
       --go_out=$PWD/internal/block --go_opt=paths=source_relative $PWD/protos/block.proto \
       --go-grpc_out=require_unimplemented_servers=false:$PWD/internal/block  --go-grpc_opt=paths=source_relative $PWD/protos/block.proto