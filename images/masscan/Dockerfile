FROM ubuntu:18.04

RUN apt-get update

RUN apt-get --assume-yes install git make gcc libpcap-dev  inetutils-ping

ADD src /src

RUN cd /src && make && make install

ENTRYPOINT ["/usr/bin/masscan"]
