FROM golang:alpine

WORKDIR /build

ADD src /build

ENV GOPROXY=https://proxy.golang.com.cn,direct
ENV GOPRIVATE=git.mycompany.com,github.com/my/private

RUN echo "https://mirrors.tuna.tsinghua.edu.cn/alpine/v3.15/main" > /etc/apk/repositories && echo "https://mirrors.tuna.tsinghua.edu.cn/alpine/v3.15/community" >> /etc/apk/repositories

RUN apk update && apk add build-base

RUN  mkdir /app &&  go build -trimpath -ldflags="-s -w"  -o /app/rotateproxy cmd/rotateproxy/main.go

WORKDIR /app


ENTRYPOINT ["/app/rotateproxy"]
