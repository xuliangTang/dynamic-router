FROM golang:1.19-alpine as builder
RUN mkdir /src
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
ADD . /src
WORKDIR /src
RUN GOPROXY=https://goproxy.cn go build -o dynamic-router main.go  && chmod +x dynamic-router


FROM alpine:3.12
RUN mkdir /app
WORKDIR /app
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk add tzdata
ENV TZ=Asia/Shanghai
ENV ZONEINFO=/app/zoneinfo.zip

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /app

COPY --from=builder /src/dynamic-router /app
COPY --from=builder /src/config.yaml /app

ENTRYPOINT  ["./dynamic-router"]