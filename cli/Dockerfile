FROM golang:1.16-alpine

WORKDIR /cliapp

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY *.go ./

RUN go build -o /cli

EXPOSE 10002

CMD ["/cli"]