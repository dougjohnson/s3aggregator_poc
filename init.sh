#!/bin/bash
echo "-- building s3aggregator container"
docker images | grep s3aggregator | awk '{print $3'} | xargs docker rmi
docker build -t s3aggregator .
echo "\n"
echo "-- success!"
echo "-- build.sh can now be used to compile a linux executable in this dir"
echo "-- test.sh can now be used to run the tests in this dir"
echo "-- no need to install golang on your local machine."
