FROM golang:alpine as builder

RUN apk update \
  && apk add --no-cache git curl

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY static static
COPY main.go .

RUN GOOS=linux GOARCH=amd64 go build -o main

FROM alpine:3.13.5
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/static static

ENV PORT=${PORT}
ENTRYPOINT ["/app/main"]