DIR=/go/src/github.com/dmcsorley/goblin
LOGOPT=--log-opt max-size=10m --log-opt max-file=5
SOCKV=-v /var/run/docker.sock:/var/run/docker.sock
LOGSPOUTIGNORE=-e LOGSPOUT=ignore
EXAMPLEIMAGE=dmcsorley/goblin-example

.PHONY: build fmt inc deps goblin example runlogstash runlogspout runexample runall

build: deps goblin

fmt:
	find . -name "*.go" | xargs dirname | sort | uniq | xargs go fmt

inc:
	docker run -it --rm -w $(DIR) -v $$PWD:$(DIR) dmcsorley/goblin:deps bash

deps:
	docker build --pull=true --no-cache -t dmcsorley/goblin:deps -f Dockerfile.deps .

buildfromdeps:
	docker run -it --rm -w $(DIR) -v $$PWD:$(DIR) dmcsorley/goblin:deps bash -c "go install -v && cp /go/bin/goblin ./bin/"

goblinimage:
	docker build --pull=true --no-cache -t dmcsorley/goblin .
	
goblin: buildfromdeps goblinimage

example:
	cd example && docker build --no-cache -t $(EXAMPLEIMAGE) .

runlogstash:
	docker run -d --name logstash $(LOGOPT) -v $$PWD/example/logstash.conf:/etc/logstash.conf $(LOGSPOUTIGNORE) -p 5000:5000 logstash -f /etc/logstash.conf

runlogspout:
	docker run -d --name logspout $(LOGOPT) $(SOCKV) $(LOGSPOUTIGNORE) gliderlabs/logspout \
	syslog://$(shell docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' logstash):5000?filter.name=goblin-*

runexample:
	docker run -d $(LOGOPT) $(SOCKV) -e IMAGE=$(EXAMPLEIMAGE) --name goblin-example -p 8080:80 $(EXAMPLEIMAGE)

runall:
	make runlogstash
	sleep 5
	make runlogspout
	sleep 5
	make runexample
