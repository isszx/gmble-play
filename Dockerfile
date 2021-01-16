FROM golang:1.15.6-alpine

RUN apk update
RUN apk add --no-cache --update git build-base opus ffmpeg openal-soft python3 py3-pip
RUN python3 -m ensurepip
RUN pip3 install --no-cache --upgrade pip setuptools
RUN pip3 install youtube-dl

ADD . /go/src/gmble-play
WORKDIR /go/src/gmble-play
RUN go install