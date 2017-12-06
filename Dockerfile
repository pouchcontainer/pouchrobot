FROM ubuntu:16.04

# install wget to download golang source code
# install git
RUN apt-get update \
    && apt-get install -y \
    wget \ 
    git \
    make \
    gcc \
    vim \
    tree \
    && apt-get clean

# set go version this image use
ENV GO_VERSION=1.9.1
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

EXPOSE 6789

COPY . /go/src/github.com/allencloud/automan

RUN go get github.com/allencloud/automan

# try to skip StrictHostKeyChecking when executing git
RUN echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config

RUN mkdir -p /go/src/github.com/alibaba \
    && cd /go/src/github.com/alibaba \
    && git clone https://github.com/pouchrobot/pouch.git \
    && cd pouch \
    && git remote remove origin \
    && git remote add origin git@github.com:pouchrobot/pouch.git \
    && git remote add upstream https://github.com/alibaba/pouch.git \
    && git config user.name "pouchrobot" \ 
    && git config user.email "pouch-dev@alibaba-inc.com"

WORKDIR /go/src/github.com/alibaba/pouch
