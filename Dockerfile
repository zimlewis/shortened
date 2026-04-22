# Build
FROM golang:1.26.2-alpine AS build-stage

RUN apk add --no-cache gcc g++ make libwebp-dev
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

WORKDIR /app
COPY go.mod go.sum .
RUN go mod download

COPY ./ .
RUN go build -o server .

# Run
FROM alpine:3.22.4 AS run

COPY --from=build-stage /app/server .

CMD ["./server"]
