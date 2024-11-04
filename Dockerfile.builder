FROM golang:1.23-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apk update && apk add --no-cache docker

RUN go build -o myapp ./cmd/builder/main.go

CMD ["./myapp"]

# # Start a new stage from scratch
# FROM alpine:latest  

# # Copy the Pre-built binary file from the previous stage
# COPY --from=builder /app/myapp .

# # Command to run the executable
# CMD ["./myapp"]
