# Build stage
FROM golang:1.22 AS builder
ENV CGO_ENABLED 0

WORKDIR /app

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all dependencies
RUN go mod tidy

# Build the Go app
RUN go build -ldflags='-s -w' -o bin/app main.go


# Check if the binary was created
RUN if [ ! -f bin/app ]; then echo "Go build failed"; exit 1; fi

# Stage 2: Create a minimal runtime image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/app /bin/app
COPY --from=builder /app/config.yaml  ./config.yaml


# Check if the binary was copied
RUN if [ ! -f /bin/app ]; then echo "Binary not copied"; exit 1; fi

CMD ["/bin/app", "-c", "/etc/config.yaml"]