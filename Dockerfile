FROM golang:1.25.0-alpine3.21
WORKDIR /app
COPY . .
# change this to config-uat if you want to use config uat config
RUN export CONFIG_FILE=config-prod
RUN go mod tidy
RUN go build
ENTRYPOINT ["/app/app"]