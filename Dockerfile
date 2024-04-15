FROM golang:1.20-buster as builder

# Build args.
ARG SHASUM 

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -ldflags "-X main.SHASUM=$SHASUM" -o bin/api cmd/api/main.go

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM debian:buster-slim

# Set up dpkg.
RUN dpkg --configure -a

# Install basic dependencies.
RUN set -x && apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates wget unzip && \
    rm -rf /var/lib/apt/lists/*


# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/bin/api /app/bin/api

# Run the web service on container startup.
ENTRYPOINT [ "/app/bin/api" ]
