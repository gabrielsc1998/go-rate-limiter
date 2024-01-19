FROM golang:1.21.1-alpine

ENV USER=root

RUN go install github.com/cosmtrek/air@latest

WORKDIR /go/src
COPY . .

# CMD ["tail", "-f", "/dev/null"]

# Run the application
ENTRYPOINT [ "sh", "./.docker/entrypoint.sh" ]
