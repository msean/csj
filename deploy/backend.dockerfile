FROM alpine:latest

# 安装时区数据
RUN apk add --no-cache tzdata

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

RUN mkdir -p /etc/caishuji/log
COPY backend_bin /usr/local/bin/server
RUN chmod +x /usr/local/bin/server

WORKDIR /

ENV ORCH_CONF=/etc/caishuji/backend_conf.yaml

STOPSIGNAL SIGRTMIN+3

ENTRYPOINT server -c $ORCH_CONF
