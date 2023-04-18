FROM golang:1.20 as dev
EXPOSE 40000
WORKDIR /app

ARG AIR_VERSION=v1.42.0

RUN DEBIAN_FRONTEND=noninteractive apt update && apt upgrade -y \
    && apt install -y git libvips-dev \
    make openssh-client

RUN go install github.com/cosmtrek/air@${AIR_VERSION} \
    && go install github.com/go-delve/delve/cmd/dlv@latest

CMD air

FROM golang:1.20 as testing
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN DEBIAN_FRONTEND=noninteractive apt update && apt upgrade -y \
    && apt install -y libvips-dev && \
    go mod download

CMD go test -vet=off -v ./test/...


FROM golang:1.20 AS build
LABEL stage="gobuilder"
WORKDIR /build

ENV CGO_ENABLED 1
ENV GOOS linux

RUN DEBIAN_FRONTEND=noninteractive apt update && \
    apt install -y libvips-dev

COPY . .
RUN go build -trimpath -ldflags="-s -w" -o /app/server cmd/api/main.go


FROM golang:1.20 as prod
WORKDIR /app
ENV TZ America/New_York

RUN DEBIAN_FRONTEND=noninteractive apt update && \
    apt install -y libvips ca-certificates

COPY --from=build /app/server /app/server

CMD ["./server"]


FROM prod as test
