#! /bin/bash
docker build . -t hello-world-server
docker run --rm -d -p4000:4000 hello-world-server

sleep 3
curl localhost:4000