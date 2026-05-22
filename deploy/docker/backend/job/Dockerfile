FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app
COPY job-rpc .
RUN mkdir -p /app/etc && \
    chmod +x /app/job-rpc

ENTRYPOINT ["./job-rpc", "-f", "/app/etc/job.yaml"]
