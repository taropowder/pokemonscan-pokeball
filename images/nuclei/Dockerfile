FROM ubuntu:18.04

ARG version=2.8.9
WORKDIR /app

ADD https://github.com/projectdiscovery/nuclei/releases/download/v${version}/nuclei_${version}_linux_amd64.zip /app
ADD https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb /app

#ADD nuclei_2.8.8_linux_amd64.zip /app
#ADD google-chrome-stable_current_amd64.deb /app

RUN apt-get update

RUN apt install -y unzip curl net-tools iputils-ping && \
    apt install -y ./google-chrome-stable_current_amd64.deb

RUN mkdir /app/res

RUN unzip nuclei_${version}_linux_amd64.zip

RUN  /app/nuclei -update

RUN  /app/nuclei -ut

#delete some useless plugins 

# disable ssl
RUN rm -rf /root/nuclei-templates/ssl
# disable harmless http detection
RUN rm -f /root/nuclei-templates/misconfiguration/http-missing-security-headers.yaml
RUN rm -f /root/nuclei-templates/misconfiguration/httponly-cookie-detect.yaml



RUN  /app/nuclei -headless

ENTRYPOINT ["/app/nuclei"]
#RUN mkdir /app

# docker run  -v  --rm pokemon:plugin_nuclei   -es info,low  -json -o res.json  -u http://host.docker.internal:8161
