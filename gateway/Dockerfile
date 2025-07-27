# All paths here will be relative to the root of the project, as that's where you'll
# need to run the docker build command from
FROM golang:1.23.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY gateway /app/gateway
COPY libs /app/libs

RUN go build -o gateway-binary ./gateway/main.go

# permission issues
RUN chmod +x /app/gateway-binary

EXPOSE 8080

CMD ["./gateway-binary"]