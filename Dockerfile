FROM golang:1.23.4

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy go mod and sum files
COPY go.* ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o mqtt-cli.exe

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["start"]
