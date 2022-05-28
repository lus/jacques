# Build the binary on the golang image
FROM golang:1.18-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o jacques ./cmd/bot/main.go

# Run the binary on distroless
FROM gcr.io/distroless/base:latest
WORKDIR /root
COPY --from=build /app/jacques .
CMD ["./jacques"]
