From alpine:latest

COPY backend_bin /usr/local/bin/

WORKDIR /

ENV ORCH_CONF=/etc/caishuji/backend_conf.yaml

STOPSIGNAL SIGRTMIN+3

ENTRYPOINT server -c $ORCH_CONF
