FROM golang:1.24

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY internal ./internal
COPY docs ./docs

# Build
RUN GOOS=linux go build -o /game-api

EXPOSE 8080

# Run
CMD ["/game-api "]