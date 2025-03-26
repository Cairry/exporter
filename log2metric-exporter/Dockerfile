FROM registry.js.design/library/golang:1.23.3-alpine3.19 AS build

ARG VERSION

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /root

COPY . /root

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \
    && apk upgrade && apk add --no-cache --virtual .build-deps \
    ca-certificates upx

RUN CGO_ENABLED=0 go build -o log2metric-exporter . && \
    chmod +x log2metric-exporter

FROM registry.js.design/base/alpine-aws:3.14

COPY --from=build /root/log2metric-exporter /app/log2metric-exporter

WORKDIR /app

ENTRYPOINT ["/app/log2metric-exporter"]