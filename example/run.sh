#!/bin/bash -ex
docker run \
 -v /var/run/docker.sock:/var/run/docker.sock \
 -e IMAGE=dmcsorley/goblin-example \
 -p 8080:80 dmcsorley/goblin-example
