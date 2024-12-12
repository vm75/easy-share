# Stage 1: Build easy-share
FROM golang:alpine AS build

# Install required build tools
RUN apk add --no-cache build-base

# Set the working directory
WORKDIR /workdir

# Set environment variables for cross-compilation
ENV CGO_ENABLED=1

# Copy Go modules files first for caching
COPY server/go.mod server/go.sum /workdir/server/
RUN go mod download -C /workdir/server

# Copy the rest of the application source code
COPY server /workdir/server

# Build the application for multiple architectures
# The build command will later be specified by Buildx
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux GOARCH=${TARGETARCH} go build -C /workdir/server -ldflags="-s -w" -o /workdir/easy-share

# Stage 2: Create the final minimal image
FROM alpine:latest AS runtime

# Install required packages
RUN apk --no-cache update
RUN apk --no-cache upgrade
RUN apk --no-cache --no-progress add tini samba nfs-utils

# Copy binaries from build stage
COPY --from=build /workdir/easy-share /opt/easy-share/easy-share

# Copy the server code
COPY usr /usr
COPY server/static /opt/easy-share/static

# Create the data directory
RUN mkdir -p /data

# Set the volume
VOLUME ["/data"]

# Expose ports for easy-share, smbd, mnbd and nfs
EXPOSE 80/tcp 111 137 138 139 445 2049 20048

# Healthcheck
HEALTHCHECK --interval=60s --timeout=15s --start-period=120s \
    CMD netstat -an | grep -c ":::80 "

# Run easy-share
ENTRYPOINT [ "/opt/easy-share/easy-share" ]
