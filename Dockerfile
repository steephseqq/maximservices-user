FROM golang:1.25.1 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

FROM alpine:3.21.3
WORKDIR /app
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=builder /app/server .
RUN chown -R appuser:appgroup /app
USER appuser
EXPOSE 50052
CMD ["./server"]
