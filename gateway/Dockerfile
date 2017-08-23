FROM golang:latest
WORKDIR /go/src/github.com/wwgberlin/baby_janus/gateway
RUN cd /go/src/github.com/wwgberlin/baby_janus/gateway
CMD go test ./... && go build && /go/src/github.com/wwgberlin/baby_janus/gateway/gateway
EXPOSE 8080