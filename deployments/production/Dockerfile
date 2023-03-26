FROM golang:1.20 AS build
LABEL stage="gobuilder"
WORKDIR /build

ENV CGO_ENABLED 1
ENV GOOS linux

RUN DEBIAN_FRONTEND=noninteractive apt update && \
    apt install -y libvips-dev

COPY . .
RUN go build -trimpath -ldflags="-s -w" -o /app/server cmd/api/main.go

FROM golang:1.20
WORKDIR /app
ENV TZ America/New_York

RUN DEBIAN_FRONTEND=noninteractive apt update && \
    apt install -y libvips ca-certificates

COPY --from=build /app/server /app/server

CMD ["./server"]
