FROM ubuntu:18.04

ARG version=1.9.4
WORKDIR /app

ADD https://github.com/chaitin/xray/releases/download/${version}/xray_linux_amd64.zip /app

RUN apt-get update

RUN apt install -y unzip ibpcap0.8-dev libuv1-dev

RUN unzip xray_linux_amd64.zip && mv xray_linux_amd64 xray


RUN chmod +x /app/xray

RUN /app/xray

ENTRYPOINT ["/app/xray"]