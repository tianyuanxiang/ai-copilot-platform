FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app
COPY blade-rpc .
RUN mkdir -p /app/etc /app/blade_jobs && \
    chmod +x /app/blade-rpc

EXPOSE 9300

ENTRYPOINT ["./blade-rpc", "-f", "/app/etc/blade.yaml"]
