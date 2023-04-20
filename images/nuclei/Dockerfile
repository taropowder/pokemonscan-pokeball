FROM ubuntu:18.04

ARG version=2.9.2
WORKDIR /app

ADD https://github.com/projectdiscovery/nuclei/releases/download/v${version}/nuclei_${version}_linux_amd64.zip /app
ADD https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb /app

#ADD nuclei_${version}_linux_amd64.zip /app
#ADD google-chrome-stable_current_amd64.deb /app

RUN apt-get update

RUN apt install -y unzip curl net-tools iputils-ping git && \
    apt install -y ./google-chrome-stable_current_amd64.deb

RUN mkdir /app/res

RUN unzip nuclei_${version}_linux_amd64.zip

RUN  /app/nuclei -update

RUN  /app/nuclei -ut


RUN  /app/nuclei -headless

#delete some useless plugins 

# disable ssl
#RUN rm -rf /root/nuclei-templates/ssl
## disable harmless http detection
#RUN rm -f /root/nuclei-templates/misconfiguration/http-missing-security-headers.yaml
#RUN rm -f /root/nuclei-templates/misconfiguration/httponly-cookie-detect.yaml


#RUN  mkdir /app/nuclei-templates && /app/nuclei -ud nuclei-templates

#RUN #git clone https://github.com/projectdiscovery/nuclei-templates.git

ENTRYPOINT ["/app/nuclei"]
#RUN mkdir /app

# docker run  -v  --rm pokemon:plugin_nuclei   -es info,low  -json -o res.json  -u http://host.docker.internal:8161
#  docker build  --build-arg https_proxy=$HTTP_PROXY --build-arg http_proxy=$HTTP_PROXY   --build-arg HTTP_PROXY=$HTTP_PROXY --build-arg HTTPS_PROXY=$HTTP_PROXY  -t debug_plugin  .