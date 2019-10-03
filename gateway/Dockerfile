FROM golang:latest
WORKDIR /go/src/github.com/ronna-s/baby-janus/gateway
RUN cd /go/src/github.com/ronna-s/baby-janus/gateway
CMD go test ./... && go run /go/src/github.com/ronna-s/baby-janus/gateway/cmd/gateway/main.go
EXPOSE 8080