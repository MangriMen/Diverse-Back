FROM golang:1.20-alpine AS build-stage

LABEL stage="gobuilder"

ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build
COPY . .
RUN go build -ldflags="-s -w" -o /app/server cmd/api/main.go

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates
COPY --from=build-stage /usr/share/zoneinfo/America/New_York /user/share/zoneinfo/America/New_York
ENV TZ America/New_York

WORKDIR /app
COPY --from=build-stage /app/server /app/server
COPY --from=build-stage /build/assets /app/assets

CMD ["./server"]