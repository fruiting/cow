#!/bin/sh -e

run(){
  docker run -p 8080:8080 \
    -d \
    --env-file ./ci/config/dev.env \
    --name cow-backend \
    --network cow-net \
    romaspirin/cow-backend
}

unit(){
  echo "run unit tests..."
  go test ./...
}

unit_race() {
  echo "run unit tests with race test"
  go test -race ./...
}

command="$1"
if [ -z "$command" ]
then
 using
 exit 0;
else
 $command $@
fi