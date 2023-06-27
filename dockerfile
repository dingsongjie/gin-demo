############################
# STEP 1 build executable binary
############################
FROM golang:1.20.5 AS builder

# Install dependencies
WORKDIR /usr/src/app
COPY . .

# Fetch dependencies.
# Using go get.
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

ENV PORT 8080
ENV GIN_MODE release
EXPOSE 8080

# Run the Go Gin binary.
ENTRYPOINT ["/usr/local/bin/app"]
