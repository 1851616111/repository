FROM golang:1.5.1

ENV TIME_ZONE=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TIME_ZONE /etc/localtime && echo $TIME_ZONE > /etc/timezone

WORKDIR /go/src/github.com/asiainfoLDP/datahub_repository
ADD . /go/src/github.com/asiainfoLDP/datahub_repository

EXPOSE 8089

ENV SERVICE_NAME=datahub_repository

RUN GO15VENDOREXPERIMENT=1 go build

ENTRYPOINT ["/go/src/github.com/asiainfoLDP/datahub_repository/datahub_repositor"]


