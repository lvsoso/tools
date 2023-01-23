#!/bin/bash

cleanup() {
    echo "do something for graceful stop..."
    exit
}

trap cleanup SIGTERM

i=1
while true;
do
  echo "running for $i times"
  ((i++))
  sleep 1
done