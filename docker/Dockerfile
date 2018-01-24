FROM alpine
MAINTAINER Thomas Rampelberg <thomasr@vapor.io>

ENV DOCKER_OS=linux
ENV DOCKER_ARCH=amd64
ENV SYNCTHING_VERSION=0.14.43

RUN wget -q https://github.com/syncthing/syncthing/releases/download/v${SYNCTHING_VERSION}/syncthing-linux-amd64-v${SYNCTHING_VERSION}.tar.gz \
  && tar -xzf syncthing-linux-amd64-v${SYNCTHING_VERSION}.tar.gz syncthing-linux-amd64-v${SYNCTHING_VERSION}/syncthing \
  && mv syncthing-linux-amd64-v${SYNCTHING_VERSION}/syncthing /bin/

COPY bin/radar_${DOCKER_OS}_${DOCKER_ARCH} /radar
