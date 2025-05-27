FROM golang:1.23.5-alpine3.21
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build
ENTRYPOINT ["/app/app"]