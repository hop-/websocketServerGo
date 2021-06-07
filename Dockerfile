FROM golang:1.16.4
LABEL maintainer=priotix


WORKDIR /var/www/wss
COPY . .

RUN go build 

CMD ["./wss"]
