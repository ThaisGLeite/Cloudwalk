# STEP 1: Build executable binary
FROM golang:1.20-alpine3.18 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum from host to container
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Create the /tmp directory to store the games.log while running
RUN mkdir -p /tmp

# Copy everything from the current directory of the host to the Working Directory in the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o quake-parser ./cmd

# STEP 2: Build a small image
FROM scratch

# Copy our static and executable
COPY --from=builder /app/quake-parser /quake-parser
COPY --from=builder /app/static /static
COPY --from=builder /tmp /tmp

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the quake-parser binary
ENTRYPOINT ["/quake-parser"]
