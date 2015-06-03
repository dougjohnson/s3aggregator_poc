#!/bin/bash 
set -e

docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp s3aggregator sh -c "go build -v"
mv myapp ./s3aggregator
