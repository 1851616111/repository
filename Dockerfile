FROM golang:1.5.1

WORKDIR /go/src/github.com/asiainfoLDP/datahub_repository
ADD . /go/src/github.com/asiainfoLDP/datahub_repository

RUN go get github.com/tools/godep

RUN godep go build

EXPOSE 8089

ENV SERVICE_NAME=datahub_repository

ENTRYPOINT ["/go/src/github.com/asiainfoLDP/datahub_repository/datahub_repository"]


