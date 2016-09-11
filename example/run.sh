#!/bin/bash -ex
docker run \
 -v /var/run/docker.sock:/var/run/docker.sock \
 -p 8080:80 dmcsorley/simpleci-example
