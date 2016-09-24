build "goblin" {
  step git-clone {
    url = "https://github.com/dmcsorley/goblin"
  }
  step docker-run {
    image = "golang"
    dir = "/go/src/github.com/dmcsorley/goblin"
    cmd = "go get -v -d && go install -v && cp /go/bin/goblin ./bin/"
  }
  step docker-build {
    image = "dmcsorley/goblin"
  }
}

build captainhook {
  step git-clone {
    url = "https://github.com/dmcsorley/captainhook"
  }

  step docker-build {
    image = "dmcsorley/captainhook"
  }
}