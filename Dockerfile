FROM debian:bullseye-slim
RUN rm /bin/sh && ln -s /bin/bash /bin/sh

RUN apt-get update && \
	apt-get install -y \
	build-essential make nano git wget curl bison

WORKDIR /root
RUN curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer > ./gvm-installer.sh &&\
	chmod +x ./gvm-installer.sh &&\
	./gvm-installer.sh &&\
	source /root/.gvm/scripts/gvm &&\
	gvm install go1.14.4 -B

WORKDIR /root/solace
COPY solclient/ /root/solace/solclient/
COPY solwrap/ /root/solace/solwrap/
COPY src/ /root/solace/src/
COPY Makefile /root/solace/Makefile
COPY build.sh /root/solace/build.sh

RUN source /root/.gvm/scripts/gvm &&\
	gvm use go1.14.4 --default &&\
	./build.sh