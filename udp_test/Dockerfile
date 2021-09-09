FROM golang:1.16-alpine

WORKDIR /netapp

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /hecnet

EXPOSE 8080

CMD ["/hecnet"]