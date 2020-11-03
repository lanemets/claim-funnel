#!/bin/bash

cd proto || return

protoc --go_out=../external/profile/gen --go_opt=paths=source_relative \
  --go-grpc_out=../external/profile/gen --go-grpc_opt=paths=source_relative profile.proto

echo "Profile has been generated"

protoc --go_out=../external/claim/gen --go_opt=paths=source_relative \
  --go-grpc_out=../external/claim/gen --go-grpc_opt=paths=source_relative claim.proto

echo "Claim has been generated"
