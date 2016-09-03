# simpleci

I would like a CI system that is lightweight but powerful.
I think there are a minimal amount of functions that such a
system would implement to do everything I need. Among them are:

- receiving a signal that a commit has been pushed
- cloning a code repository
- running a docker build
- pushing to a docker repository
- notifying stakeholders of success or failure

## Things I want

- I would like the CI system I run to be instantiable.
- I would like to version control the config for my CI instance.
- I would like my CI system to scale across available resources.
- I want to not manage tools/executables/libraries on build slaves.
