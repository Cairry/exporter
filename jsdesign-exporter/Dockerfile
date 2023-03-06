FROM golang:1.18 AS build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.io"
WORKDIR /root
COPY . /root
RUN go mod tidy && \
    go build -o jsdesign-exporter ./jsdesign-exporter.go && \
    chmod 777 jsdesign-exporter

FROM busybox:glibc
RUN mkdir -p /app/exporter
COPY --from=build /root/jsdesign-exporter /app/exporter

WORKDIR /app/exporter
ENTRYPOINT ["/app/exporter/jsdesign-exporter"]