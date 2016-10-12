PACKAGE=github.com/dmcsorley/goblin
DIR=/go/src/$(PACKAGE)

.PHONY: fromdeps test fmt inc deps image goblin

fromdeps:
	docker run -it --rm -w $(DIR) -v $$PWD:$(DIR) dmcsorley/goblin:deps bash -c "go install -v && cp /go/bin/goblin ./bin/"

test:
	go test -v $(PACKAGE) $(PACKAGE)/cibuild $(PACKAGE)/config

fmt:
	go fmt $(PACKAGE) $(PACKAGE)/cibuild $(PACKAGE)/command $(PACKAGE)/config $(PACKAGE)/gobdocker

inc:
	docker run -it --rm -w $(DIR) -v $$PWD:$(DIR) dmcsorley/goblin:deps bash

deps:
	docker build --pull=true --no-cache -t dmcsorley/goblin:deps -f Dockerfile.deps .

image:
	docker build --pull=true --no-cache -t dmcsorley/goblin .

goblin: deps fromdeps image
