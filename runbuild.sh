DIR=/go/src/github.com/dmcsorley/goblin
docker run -it --rm \
 -w $DIR -v $PWD:$DIR \
 golang \
 bash -c "go get -v -d && go install -v && cp /go/bin/goblin ./bin/"
