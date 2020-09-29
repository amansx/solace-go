FROM debian:bullseye-slim
RUN rm /bin/sh && ln -s /bin/bash /bin/sh

RUN apt-get update && \
	apt-get install -y \
	build-essential make nano git wget curl bison \
	gcc-mingw-w64-x86-64 g++-mingw-w64-x86-64 wine64 \
	gcc-multilib mingw-w64-tools

WORKDIR /root
RUN curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer > ./gvm-installer.sh &&\
	chmod +x ./gvm-installer.sh &&\
	./gvm-installer.sh &&\
	source /root/.gvm/scripts/gvm &&\
	gvm install go1.14.4 -B

WORKDIR /root/solace

COPY *.c  /root/solace/
COPY *.go  /root/solace/
COPY *.mod /root/solace/
COPY makefile /root/solace/
COPY makefile.* /root/solace/
COPY build.sh /root/solace/build.sh
COPY includes/ /root/solace/includes/
COPY examples/ /root/solace/examples/

RUN source /root/.gvm/scripts/gvm && gvm use go1.14.4 --default

CMD source /root/.gvm/scripts/gvm &&\
	gvm use go1.14.4 --default &&\
	./build.sh