FROM golang:1.20
EXPOSE 40000
WORKDIR /app

RUN DEBIAN_FRONTEND=noninteractive apt update && apt upgrade -y \
    && apt install -y git libvips-dev \
    make openssh-client

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

RUN go install github.com/go-delve/delve/cmd/dlv@latest

CMD air
