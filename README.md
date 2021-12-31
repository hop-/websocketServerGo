# wss

wss is a websocket service which allaws client to authorize, get notification list and get updates if there is new notification

## Table of Contents
  1. [Requirements](#requirements)
  2. [Deployment](#deployment)
  3. [Build](#build)
  4. [Running The Service](#running-the-service)
  5. [Kafka Messaging](#kafka-messaging)
  6. [Websocket API](#websocket-api)
  7. [Backend API Requests](#backend-api-requests)

## Requirements

- Docker v18.09

- Golang v1.16 or higher
https://golang.org/doc/install?download=go1.16.linux-amd64.tar.gz

## Deployment

2. Get Sources

```sh
$ cd /home/USER/code/
$ git clone git@github.com:hop-/wss.git
```

3. Setup Environment Variables

```sh
$ cp .env.dist .env
```

  in .env replace templates with real values

```
HOST_ENV={{env}}
NODE_ENV={{env}}
KAFKA_BROKERS={{kafka-broker1:port1,kafka-broker2:port2}}
```

## Build

### Build Without Docker

```sh
$ cd path/to/wss
$ go build
```
### Build With Docker

```sh
$ cd path/to/flash
$ docker build
```

## Running The Service

### Executable

```sh
$ ./wss
```

### Docker

 Up docker container

## Kafka Messaging

### Consumer

This service consumes kafka topic for a new notifications and send it to the users.
As a kafka message consumer this service has the following configurations

#### Topic
`wss`

#### Message
```json
{
  "type": "message type",
  "document": {
    "type related payload"
  }
}
```
## Websocket API

Websocket connection is established by client request.
For each connection the service creates a separate session with unique session ID.
Communication between the client and the server is done using json objects.
You can have two types of json objects.
### Request object
Client can send request objects.
Each reqeust object is an command, and has unique reqeust id.
The options object can be different depending to the command type.
```json
{
  "rid": 1,
  "command": "command name",
  "options": {
    "options content"
  }
}
```
### Response object
Server sends response object for each client reqeust.
Each reqeust will receive response with the same reqeust id.
There is another type of response which reqeust id is equal to -1.
These responses are push notifications which client can received if subscribe to some group of notifications.
```json
{
  "rid": 1,
  "event": "event name",
  "status": "response status",
  "payload": {
    "paylod content"
  }
}
```
## Backend API Requests
This service uses http requesting mechanism to communicate with backend services (e.g.  api-user).
Service sends request to backend apis for some commands received from client.
