FROM golang:1.20-alpine3.17 AS build
LABEL stage="gobuilder"
WORKDIR /build

ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk update --no-cache && apk add --no-cache tzdata

COPY . .
RUN go build -trimpath -ldflags="-s -w" -o /app/server cmd/api/main.go

FROM alpine:3.17
WORKDIR /app
ENV TZ America/New_York

RUN apk update --no-cache && apk add --no-cache ca-certificates

COPY --from=build /usr/share/zoneinfo/America/New_York /user/share/zoneinfo/America/New_York
COPY --from=build /app/server /app/server

CMD ["./server"]