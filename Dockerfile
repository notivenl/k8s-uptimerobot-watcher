
FROM golang:1.20.2-alpine3.16

WORKDIR /app

RUN apk update
RUN apk add --no-cache git make

COPY . .

RUN make setup
RUN make build

CMD ["/app/bin/uptime"]