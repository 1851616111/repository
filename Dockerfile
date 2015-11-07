FROM golang:1.5.1

ADD . ./go/src/datahub_repository
WORKDIR ./go/src/datahub_repository
run curl -s https://raw.githubusercontent.com/pote/gpm/v1.3.2/bin/gpm | bash && \
go build

CMD ["./datahub_repository"]
