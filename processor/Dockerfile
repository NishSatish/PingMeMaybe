# All paths here will be relative to the root of the project, as that's where you'll
# need to run the docker build command from
FROM golang:1.23.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY processor /app/processor
COPY libs /app/libs

RUN go build -o processor-binary ./processor/main.go

# permission issues
RUN chmod +x /app/processor-binary

EXPOSE 8080

CMD ["./processor-binary"]