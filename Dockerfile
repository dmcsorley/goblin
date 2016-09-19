FROM buildpack-deps:xenial-scm

EXPOSE 80

CMD goblin

RUN mkdir -p /usr/local/etc/goblin/

WORKDIR /usr/local/etc/goblin/

RUN cd /tmp && \
  curl -fsL "https://get.docker.com/builds/Linux/x86_64/docker-1.12.1.tgz" -o docker.tgz && \
  echo "05ceec7fd937e1416e5dce12b0b6e1c655907d349d52574319a1e875077ccb79 *docker.tgz" | sha256sum -c && \
  tar xf docker.tgz docker/docker && \
  mv docker/docker /usr/local/bin && \
  rm -rf docker docker.tgz

ADD bin/goblin /usr/local/bin/
