#!/bin/bash -ex

# https://www.elastic.co/guide/en/logstash/current/configuration-file-structure.html
docker run -d \
 --name logstash \
 --log-opt max-size=10m --log-opt max-file=5 \
 -v $PWD/logstash.conf:/etc/logstash.conf \
 -e LOGSPOUT=ignore \
 -p 5000:5000 \
 logstash \
 -f /etc/logstash.conf

sleep 2

LOGSTASH_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' logstash)

# https://github.com/gliderlabs/logspout
docker run -d \
 --name logspout \
 --log-opt max-size=10m --log-opt max-file=5 \
 -v /var/run/docker.sock:/var/run/docker.sock \
 -e LOGSPOUT=ignore \
 gliderlabs/logspout \
 syslog://${LOGSTASH_IP}:5000?filter.name=goblin-*

sleep 5

docker build -t dmcsorley/goblin-example . && \

docker run -d \
 --log-opt max-size=10m --log-opt max-file=5 \
 -v /var/run/docker.sock:/var/run/docker.sock \
 -e IMAGE=dmcsorley/goblin-example \
 --name goblin-example \
 -p 8080:80 dmcsorley/goblin-example
