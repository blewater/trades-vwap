FROM golang:1.17-buster

RUN mkdir -p /go/src/vwap/build
WORKDIR /go/src/vwap
ADD . ./
RUN go build -o build/vwap main.go && \
    ./build/vwap -h

ENTRYPOINT build/vwap

# A more appropriate running image could be added at this point i.e Alpine,
# ubuntu:18.04 and copy the artifacts from buster
