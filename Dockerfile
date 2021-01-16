FROM golang:1.15.6-alpine

RUN apk update
RUN apk add --no-cache --update git build-base opus youtube-dl ffmpeg openal-soft
