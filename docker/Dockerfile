FROM openjdk:8-alpine
MAINTAINER Thomas Rampelberg <thomasr@vapor.io>

RUN mkdir -p /var/run/watchman \
  && chmod a+w /var/run/watchman

RUN apk add --no-cache --update \
    bash \
    libc6-compat \
  && apk add --update \
    --repository http://dl-3.alpinelinux.org/alpine/edge/testing/ \
    --allow-untrusted \
    watchman

RUN mkdir -p mirror \
  && wget -O /mirror/mirror-all.jar http://repo.joist.ws/mirror-all.jar
COPY mirror.sh /mirror/

COPY bin/radar /radar
