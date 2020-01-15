FROM golang:1.13-alpine AS build
WORKDIR /go/src/github.com/utilitywarehouse/kube-namespace-annotations-exporter
COPY . /go/src/github.com/utilitywarehouse/kube-namespace-annotations-exporter
ENV CGO_ENABLED 0
RUN apk --no-cache add git &&\
  go get -t ./... &&\
  go test ./... &&\
  go build -o /kube-namespace-annotations-exporter .

ENTRYPOINT [ "/kube-namespace-annotations-exporter" ]
