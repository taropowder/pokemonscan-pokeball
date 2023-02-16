FROM golang:alpine

WORKDIR /build

ADD src /build

ENV GOPROXY=https://proxy.golang.com.cn,direct
ENV GOPRIVATE=git.mycompany.com,github.com/my/private

# 有bug 临时修复
RUN go get go.mongodb.org/mongo-driver
RUN go get go.mongodb.org/mongo-driver/x/mongo/driver/topology@v1.11.1
RUN go mod download
RUN go build -o /app/enscan -buildvcs=false  .

WORKDIR /app

RUN /app/enscan -v

RUN mkdir res

ENTRYPOINT ["/app/enscan"]
