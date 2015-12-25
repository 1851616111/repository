FROM golang:1.5.1

WORKDIR /go/src/github.com/asiainfoLDP
ADD . /go/src/github.com/asiainfoLDP/
RUN cd /go/src/github.com/asiainfoLDP/datahub_repository && ls
RUN cd /go/src/github.com/asiainfoLDP/datahub_repository && ls
RUN cd /go/src/github.com/asiainfoLDP/datahub_repository/mq && ls
RUN go get github.com/tools/godep
RUN godep restore
RUN godep go install

EXPOSE 8089

ENV SERVICE_NAME=datahub_repository

ENTRYPOINT ["/go/bin/datahub_repository"]
