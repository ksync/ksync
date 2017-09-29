FROM alpine

ENV HOME /tmp

COPY radar /radar

EXPOSE 40321

CMD ["/radar",  "--log-level=debug", "serve"]
