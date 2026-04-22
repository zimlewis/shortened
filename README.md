# Shortened

An overengineered link shortened to practice kafka, goroutin and channel

## Description

This project is created for me to learn more about kafka, channel and goroutin. Every time a link get cliked, the handler of the http will first find the full url of the saved link, if it success then it will fire a signal to the channel which will send it to a goroutin that run kafka message publisher and then write it to kafka server, the consumer will then take the message, then execute a database update that will increase the counter of the link

## Getting Started

If you want to run it for some reason

### Dependencies

* Having [Docker](https://docs.docker.com/engine/install/) installed

### Installing & run

* Clone this repo:
```
git clone https://github.com/zimlewis/shortened.git
```
* Run the command:
```
docker compose up
```

### Endpoint

| Method | Endpoint | Description | Request Body |
| --- | --- | --- | --- |
| `POST` | `/` | Add a new shortened link | `{"shortened": "slug", "full": "https://example.com"}` |
| `GET` | `/{id}` | Redirect to the full URL | N/A |
| `GET` | `/{id}/count` | Get the click count for a link | N/A |
| `PUT` | `/{id}` | Update the full URL of an existing link | `{"full": "https://new-url.com"}` |
| `DELETE` | `/{id}` | Delete a shortened link | N/A |

#### Examples

**Create a link**
```bash
curl -X POST http://localhost:8080/ \
     -H "Content-Type: application/json" \
     -d '{"shortened": "ggl", "full": "https://google.com"}'
```

**Get click count**
```bash
curl http://localhost:8080/ggl/count
```




## Help

See compose.yaml to setup your own environment

## Environment variable

* Badger DB location
```
BADGER_DIR
```
* Kafka broker for the reader and writer
```
KAFKA_BROKER
```
* Port the http server will run on
```
PORT
```
