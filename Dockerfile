FROM golang:1.25.0-alpine3.21
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build
ENTRYPOINT ["/app/app"]