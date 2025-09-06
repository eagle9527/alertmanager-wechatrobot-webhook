FROM alpine:latest

# 设置时区
ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata ca-certificates \
    && cp /usr/share/zoneinfo/$TZ /etc/localtime 

# 拷贝可执行文件
COPY alertmanager-wechatbot-webhook /alertmanager-wechatbot-webhook
RUN chmod +x /alertmanager-wechatbot-webhook

# 启动
ENTRYPOINT ["/alertmanager-wechatbot-webhook"]
