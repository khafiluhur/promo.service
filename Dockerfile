# syntax=docker/dockerfile:1.3
# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

# Use the secret to access the SSH key
RUN --mount=type=secret,id=ssh_key \
    mkdir -p /root/.ssh && \
    cat /run/secrets/ssh_key > /root/.ssh/id_ed25519 && \
    chmod 600 /root/.ssh/id_ed25519

# Install Git
RUN apk update && apk upgrade --no-cache && apk add --no-cache git coreutils openssh-client

RUN git config --global url."git@gitlab.com:".insteadOf https://gitlab.com/ \
    && git config --global url."git@github.com:".insteadOf https://github.com/ \
    && ssh-keyscan gitlab.com >> ~/.ssh/known_hosts \
    && ssh-keyscan github.com >> ~/.ssh/known_hosts

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go env -w GOPRIVATE=github.com/Golden-Rama-Digital && go mod download

# Copy the source code to the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /my-app .

# Stage 2: Final minimal runtime image
FROM alpine:latest

# Install only the minimal runtime dependencies (if any required, e.g., libc)
RUN apk --no-cache add ca-certificates tzdata &&  \
    ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/entrypoint.sh /entrypoint.sh
COPY --from=builder /my-app /my-app

# Set executable permission (optional since Go builds are already executable)
RUN chmod +x /entrypoint.sh

# Set the command to run the Go application
CMD ["/bin/sh", "/entrypoint.sh"]