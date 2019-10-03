FROM golang:latest
WORKDIR /go/src/github.com/ronna-s/baby-janus/server
ADD . /go/src/github.com/ronna-s/baby-janus/server
CMD go install github.com/ronna-s/baby-janus/server \
&& chmod +x /go/src/github.com/ronna-s/baby-janus/wait-for-it.sh \
&& /go/src/github.com/ronna-s/baby_janus/wait-for-it.sh baby_janus_gateway:8080/routes \
-- /go/bin/server