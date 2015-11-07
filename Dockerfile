FROM golang:1.5.1

WORKDIR /go/src/datahub_repository
ADD . /go/src/datahub_repository

RUN go get github.com/tools/godep
RUN godep restore


run curl -s https://raw.githubusercontent.com/pote/gpm/v1.3.2/bin/gpm | bash && \ godep go build

CMD["./datahub_repository"]
