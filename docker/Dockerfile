FROM golang:1.23.4-alpine3.21 as builder

ENV GOPATH=/go
WORKDIR $GOPATH/src/github.com/pegondo/starwars/service

# Install the dependencies.
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

# Build the service.
COPY ./main.go .
COPY ./internal ./internal
RUN go build -o service .

RUN ls -al


FROM alpine:latest

WORKDIR /app

# Copy the already built binary.
COPY --from=builder /go/src/github.com/pegondo/starwars/service .

# Run the binary.
CMD [ "./service" ]
