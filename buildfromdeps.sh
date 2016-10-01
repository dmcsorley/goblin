DIR=/go/src/github.com/dmcsorley/goblin

docker pull buildpack-deps:xenial-scm

docker run -it --rm \
 -w $DIR -v $PWD:$DIR \
 dmcsorley/goblin:deps \
 bash -c "go install -v && cp /go/bin/goblin ./bin/"

docker build -t dmcsorley/goblin .
