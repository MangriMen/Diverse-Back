ARG GOLANG_IMAGE_TAG="1.20-alpine3.17"

FROM golang:${GOLANG_IMAGE_TAG} as base

FROM base as dev
WORKDIR /app

EXPOSE 40000

RUN apk update && apk upgrade && \
    apk add --no-cache vips-dev gcc musl-dev && \
    pkg-config vips --cflags && \
    go install github.com/cosmtrek/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

CMD ["air"]


FROM base AS build
WORKDIR /build

ENV CGO_ENABLED 1
ENV GOOS linux

COPY go.mod .
COPY go.sum .

RUN apk update && apk upgrade && \
    apk add --no-cache vips-dev gcc musl-dev && \
    pkg-config vips --cflags && \
    go mod download

COPY . .

RUN go build -trimpath -ldflags="-s -w" -o /app/server cmd/api/main.go


FROM base as prod
WORKDIR /app
ENV TZ America/New_York

COPY --from=build /app/server /app/server

RUN apk update && apk upgrade && \
    apk add --no-cache vips-tools

CMD ["./server"]


FROM prod as test
