FROM ubuntu:16.04

# install wget to download golang source code
# install git
RUN apt-get update && \
    apt-get install -y wget git make gcc vim tree software-properties-common && \
    add-apt-repository ppa:webupd8team/java -y && \
    apt-get update && \
    echo oracle-java7-installer shared/accepted-oracle-license-v1-1 select true | /usr/bin/debconf-set-selections && \
    apt-get install -y oracle-java8-installer && \
    apt-get clean

# install swagger2markup
RUN wget --quiet -O /root/swagger2markup-cli-1.3.1.jar http://central.maven.org/maven2/io/github/swagger2markup/swagger2markup-cli/1.3.1/swagger2markup-cli-1.3.1.jar

# set go version this image use
ENV GO_VERSION=1.10.4
ENV ARCH=amd64

# install golang which version is GO_VERSION
RUN wget https://storage.googleapis.com/golang/go${GO_VERSION}.linux-${ARCH}.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-${ARCH}.tar.gz \
    && rm go${GO_VERSION}.linux-${ARCH}.tar.gz

# create GOPATH
RUN mkdir /go
WORKDIR /go
ENV GOPATH=/go

# set go binary path to local $PATH
# go binary path is /usr/local/go/bin
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# try to skip StrictHostKeyChecking when executing git
RUN echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config

RUN mkdir -p /go/src/github.com/dragonflyoss \
    && cd /go/src/github.com/dragonflyoss \
    && git clone https://github.com/pouchrobot/Dragonfly.git \
    && cd Dragonfly \
    && git remote remove origin \
    && git remote add origin git@github.com:pouchrobot/Dragonfly.git \
    && git remote add upstream https://github.com/dragonflyoss/Dragonfly.git \
    && git config user.name "pouchrobot" \ 
    && git config user.email "pouch-dev@list.alibaba-inc.com"

EXPOSE 6789

COPY . /go/src/github.com/pouchcontainer/pouchrobot

RUN cd /go/src/github.com/pouchcontainer/pouchrobot \
    && go build \
    && mv pouchrobot /usr/local/bin/pouchrobot

WORKDIR /go/src/github.com/dragonflyoss/Dragonfly
