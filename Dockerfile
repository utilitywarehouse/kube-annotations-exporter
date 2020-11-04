FROM golang:1.13-alpine AS build
WORKDIR /go/src/github.com/utilitywarehouse/kube-annotations-exporter
COPY . /go/src/github.com/utilitywarehouse/kube-annotations-exporter
ENV CGO_ENABLED 0
RUN apk --no-cache add git &&\
  go get -t ./... &&\
  go test ./... &&\
  go build -o /kube-annotations-exporter .

FROM alpine:3.10
COPY --from=build /kube-annotations-exporter /kube-annotations-exporter

ENTRYPOINT [ "/kube-annotations-exporter"]
