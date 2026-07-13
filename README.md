# Shortened

An overengineered link shortener to practice distributed streaming, concurrency, and event-driven architecture. 

## Description

A high-performance URL shortener. It uses **Go Channels** and **Goroutines** for asynchronous processing, publishing click events to **Apache Kafka** to update metrics without delaying user redirects.

## Features

- **Asynchronous Metrics:** Non-blocking processing of clicks.
- **Event-Driven:** Kafka-backed analytics.
- **High Concurrency:** Optimized for speed.

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
| `POST` | `/` | Add a new shortened link | `{"full": "https://example.com"}` |
| `GET` | `/{id}` | Redirect to the full URL | N/A |
| `GET` | `/{id}/count` | Get the click count for a link | N/A |
| `PUT` | `/{id}` | Update the full URL of an existing link | `{"full": "https://new-url.com"}` |
| `DELETE` | `/{id}` | Delete a shortened link | N/A |

#### Examples

**Create a link**
```bash
curl -X POST http://localhost:3000/ \
     -H "Content-Type: application/json" \
     -d '{"full": "https://google.com"}'
```

**Get click count**
```bash
curl http://localhost:3000/ggl/count
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
