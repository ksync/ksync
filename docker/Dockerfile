FROM alpine
MAINTAINER Thomas Rampelberg <thomasr@vapor.io>

ENV DOCKER_OS=linux
ENV DOCKER_ARCH=amd64

COPY docker/config.xml /var/syncthing/config/config.xml

RUN echo 'syncthing:x:1000:1000::/var/syncthing:/sbin/nologin' >> /etc/passwd \
    && echo 'syncthing:!::0:::::' >> /etc/shadow \
    && chown -R syncthing /var/syncthing

ENV release=

RUN apk add --no-cache --virtual .deps \
         curl \
         gnupg \
         jq \
    && apk add --no-cache \
         ca-certificates \
    && gpg --keyserver keyserver.ubuntu.com --recv-key D26E6ED000654A3E \
    && set -x \
    && mkdir /syncthing \
    && cd /syncthing \
    && release=${release:-$(curl -s https://api.github.com/repos/syncthing/syncthing/releases/latest | jq -r .tag_name )} \
    && curl -sLO https://github.com/syncthing/syncthing/releases/download/${release}/syncthing-linux-amd64-${release}.tar.gz \
    && curl -sLO https://github.com/syncthing/syncthing/releases/download/${release}/sha256sum.txt.asc \
    && gpg --verify sha256sum.txt.asc \
    && grep syncthing-linux-amd64 sha256sum.txt.asc | sha256sum \
    && tar -zxf syncthing-linux-amd64-${release}.tar.gz \
    && mv syncthing-linux-amd64-${release}/syncthing . \
    && rm -rf syncthing-linux-amd64-${release} sha256sum.txt.asc syncthing-linux-amd64-${release}.tar.gz \
    && apk del .deps

USER syncthing
ENV STNOUPGRADE=1

COPY bin/radar_${DOCKER_OS}_${DOCKER_ARCH} /radar
