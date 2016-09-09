#!/bin/bash -ex
docker run \
 -v /var/run/docker.sock:/var/run/docker.sock \
 -v $PWD/example/config.json:/go/src/app/config.json \
 -p 8080:80 dmcsorley/simpleci
