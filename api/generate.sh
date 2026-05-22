#!/bin/bash
set -e

# Change directory to the location of this script
cd "$(dirname "$0")"

# Generate the Go structs and gRPC interfaces using the source_relative path option
# so the files are placed right next to bartwarn.proto
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    bartwarn.proto
