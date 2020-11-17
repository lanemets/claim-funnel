#!/bin/bash

cd external/proto || return

protoc --go_out=../profile/gen --go_opt=paths=source_relative \
  --go-grpc_out=../profile/gen --go-grpc_opt=paths=source_relative profile.proto

echo "Profile has been generated"

protoc --go_out=../claim/gen --go_opt=paths=source_relative \
  --go-grpc_out=../claim/gen --go-grpc_opt=paths=source_relative claim.proto

echo "Claim has been generated"
