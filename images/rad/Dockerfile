FROM ubuntu:18.04

ARG version=0.4
WORKDIR /app

ADD https://github.com/chaitin/rad/releases/download/${version}/rad_linux_amd64.zip /app
ADD https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb /app

RUN apt-get update

RUN apt install -y unzip && \
    apt install -y ./google-chrome-stable_current_amd64.deb


RUN     unzip rad_linux_amd64.zip &&  mv rad_linux_amd64 rad
# && \   rm -f google-chrome-stable_current_amd64.deb

RUN mkdir /data


ADD rad_config.yml /app/rad_config.yml

ENTRYPOINT ["/app/rad"]

# docker run --rm -v /tmp/tmp_data/:/data pokemon:plugin_rad --http-proxy 172.18.173.73:7777  -t http://taropowder.cn