FROM golang:1.19

ADD src /app

WORKDIR /app

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct



RUN go mod download

RUN CGO_ENABLED=0 go build -trimpath -ldflags=" \
-w -s \
-X 'github.com/wgpsec/ENScan/common.BuiltAt=`date +'%F %T %z'`' \
-X 'github.com/wgpsec/ENScan/common.GoVersion=`go version | sed 's/go version //'`' \
-X 'github.com/wgpsec/ENScan/common.GitAuthor=`git show -s --format='format:%aN <%ae>' HEAD`' \
-X 'github.com/wgpsec/ENScan/common.GitCommit=`git log --pretty=format:"%h" -1`' \
-X 'github.com/wgpsec/ENScan/common.GitTag=`git describe --abbrev=0 --tags`' \
" \
-o /ENScan_GO/enscan .

WORKDIR /ENScan_GO

RUN /ENScan_GO/enscan -v

ENTRYPOINT ["/ENScan_GO/enscan"]