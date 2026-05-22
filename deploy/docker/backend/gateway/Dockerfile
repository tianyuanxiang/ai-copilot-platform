FROM alpine:3.19

# 安装时区数据（Go 的 time 包需要）和 CA 证书（HTTPS 请求需要）
RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

# 复制已编译好的二进制文件
COPY gateway .

# 创建必要目录
RUN mkdir -p /app/etc /app/uploads && \
    chmod +x /app/gateway

EXPOSE 6666

# 启动命令：-f 指定配置文件路径（配置文件通过 ConfigMap 挂载进来）
ENTRYPOINT ["./gateway", "-f", "/app/etc/gateway-api.yaml"]
