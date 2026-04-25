FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o gorag .

FROM alpine:latest
WORKDIR /root/

# Copy binary and static files
COPY --from=builder /app/gorag .
COPY --from=builder /app/index.html .

# Create directory for documents and database
RUN mkdir docs index_db

EXPOSE 9090

CMD ["./gorag"]
