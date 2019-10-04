FROM golang:latest
WORKDIR /go/src/github.com/ronna-s/baby-janus/app
ADD . /go/src/github.com/ronna-s/baby-janus/app
CMD go install github.com/ronna-s/baby-janus/app && /go/bin/app
EXPOSE 8080