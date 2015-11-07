FROM golang:1.5.1

WORKDIR /go/src/datahub_repository
ADD . /go/src/datahub_repository/

RUN go get github.com/tools/godep
RUN godep restore
RUN godep go install

ENV SERVICE_NAME=datahub_repository

EXPOSE 8080
ENTRYPOINT ["/go/bin/datahub_repository"]
