FROM golang:1.12.7
MAINTAINER priotix


WORKDIR /var/www/flash
COPY . .

RUN go get -d -v ./...

RUN go build .

CMD ["./flash"]
