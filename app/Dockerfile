FROM golang:latest
WORKDIR /go/src/github.com/wwgberlin/baby_janus/app
ADD . /go/src/github.com/wwgberlin/baby_janus/app
CMD go install github.com/wwgberlin/baby_janus/app \
&& chmod +x /go/src/github.com/wwgberlin/baby_janus/wait-for-it.sh \
&& /go/src/github.com/wwgberlin/baby_janus/wait-for-it.sh baby_janus_gateway:8080/routes \
-- /go/bin/app
EXPOSE 8080