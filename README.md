# goblin ci

A small CI system that uses containers to manage builds.

## Things I want

- A single config file to specify builds.
- The CI system should be versionable, buildable, and scalable like any other application.
- Take advantage of docker containers to isolate builds and make them repeatable.
- Take advantage of docker swarm to scale CI.
- Refer to CI workers as "my goblin minions".
