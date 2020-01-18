# Start from the latest golang base image
FROM golang:latest as builder

LABEL maintainer="rtkennelly1@gmail.com"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go module files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .
RUN go get -u github.com/gobuffalo/packr/v2/packr2 && cd pkg/api && packr2 && cd -

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o testlive ./cmd/testlive

######## Start a new stage from scratch #######
FROM scratch

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/testlive .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./testlive"]
