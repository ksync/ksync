FROM openjdk:8-alpine
MAINTAINER Thomas Rampelberg <thomasr@vapor.io>

RUN mkdir -p /var/run/watchman

RUN apk add --no-cache --update \
    bash \
    libc6-compat \
  && apk add --update \
    --repository http://dl-3.alpinelinux.org/alpine/edge/testing/ \
    --allow-untrusted \
    watchman

# TODO: move these local
ADD http://repo.joist.ws/mirror-all.jar /mirror/
ADD http://repo.joist.ws/mirror /mirror/

WORKDIR /mirror
CMD /bin/bash ./mirror server
