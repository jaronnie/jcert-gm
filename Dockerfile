FROM centos:7.9.2009

LABEL MAINTAINER jaron@jaronnie.com

COPY ./dist/jcert-gm_linux_amd64/jcert-gm /usr/bin/jcert-gm

RUN yum -y install bash-completion \
    && jcert-gm completion bash > /etc/bash_completion.d/jcert-gm \
    && jcert-gm init bash

