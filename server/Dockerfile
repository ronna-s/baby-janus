FROM golang:latest
WORKDIR /go/src/github.com/wwgberlin/baby_janus/server
ADD . /go/src/github.com/wwgberlin/baby_janus/server
CMD go install github.com/wwgberlin/baby_janus/server \
&& chmod +x /go/src/github.com/wwgberlin/baby_janus/wait-for-it.sh \
&& /go/src/github.com/wwgberlin/baby_janus/wait-for-it.sh baby_janus_gateway:8080/routes \
-- /go/bin/server
EXPOSE 8080