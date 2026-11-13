FROM alpine:latest

# 使用国内镜像源，避免拉取 alpine 包（尤其是 font-noto-cjk）过慢
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache tzdata font-noto-cjk

# 设置时区
ENV TZ=Asia/Shanghai

# 同步系统时区（可选但推荐）
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

RUN mkdir -p /etc/caishuji/log
COPY app_bin /usr/local/bin/server
RUN chmod +x /usr/local/bin/server

WORKDIR /

ENV ORCH_CONF=/etc/caishuji/conf.yaml

STOPSIGNAL SIGRTMIN+3

ENTRYPOINT server -c $ORCH_CONF -m true -s true
