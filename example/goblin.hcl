build "goblin" {
  step git-clone {
    url = "https://github.com/dmcsorley/goblin"
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
    dir = "captainhook"
  }
}
