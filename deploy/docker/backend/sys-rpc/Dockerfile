FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app
COPY sys-rpc .
RUN mkdir -p /app/etc /app/uploads && \
    chmod +x /app/sys-rpc

EXPOSE 9200

ENTRYPOINT ["./sys-rpc", "-f", "/app/etc/sys.yaml"]
