From alpine:latest

RUN mkdir -p /etc/caishuji/log
COPY app_bin /usr/local/bin/server
RUN chmod +x /usr/local/bin/server

WORKDIR /

ENV ORCH_CONF=/etc/caishuji/conf.yaml

STOPSIGNAL SIGRTMIN+3

ENTRYPOINT server -c $ORCH_CONF
