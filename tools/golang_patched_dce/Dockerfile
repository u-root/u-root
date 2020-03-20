FROM golang:1.13.5-alpine

ADD ad1d54f.diff /

RUN cd $(go env GOROOT) && patch -p 1 < /ad1d54f.diff

ENV GOCACHE=/go/.cache
