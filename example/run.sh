#!/bin/bash -ex
docker build -t dmcsorley/goblin-example . && \
docker run \
 -v /var/run/docker.sock:/var/run/docker.sock \
 -e IMAGE=dmcsorley/goblin-example \
 -p 8080:80 dmcsorley/goblin-example
