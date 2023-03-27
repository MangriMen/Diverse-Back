FROM golang:1.20 as dev
EXPOSE 40000
WORKDIR /app

RUN DEBIAN_FRONTEND=noninteractive apt update && apt upgrade -y \
    && apt install -y git libvips-dev \
    make openssh-client

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

RUN go install github.com/go-delve/delve/cmd/dlv@latest

CMD air


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
