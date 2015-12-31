FROM golang:1.5.1

ENV TIME_ZONE=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TIME_ZONE /etc/localtime && echo $TIME_ZONE > /etc/timezone

WORKDIR /go/src/github.com/asiainfoLDP/datahub_repository
ADD . /go/src/github.com/asiainfoLDP/datahub_repository

RUN go get github.com/tools/godep

RUN godep go build

EXPOSE 8089

ENV SERVICE_NAME=datahub_repository

ENTRYPOINT ["/go/src/github.com/asiainfoLDP/datahub_repository/datahub_repository"]


