FROM golang:1.16
LABEL maintainer=priotix

WORKDIR /wss

COPY go.mod /wss
COPY go.sum /wss
RUN ls
RUN go mod download

ADD . /wss
RUN go build 

CMD ["./wss"]
