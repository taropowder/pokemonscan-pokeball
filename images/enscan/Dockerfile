FROM golang:1.19

ADD src /app

WORKDIR /app

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct



RUN go mod download

RUN CGO_ENABLED=0 go build -trimpath -o /ENScan_GO/enscan -buildvcs=false .

WORKDIR /ENScan_GO

RUN /ENScan_GO/enscan -v

ENTRYPOINT ["/ENScan_GO/enscan"]