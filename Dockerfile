FROM golang:1.17-alpine AS compiler
WORKDIR /src/app

RUN apk add build-base

COPY go.mod ./
RUN go mod download
COPY . .

RUN go build -a -o bin/sshman *.go

ENTRYPOINT /src/app/bin/sshman