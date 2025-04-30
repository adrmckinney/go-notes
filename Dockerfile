# Use the official Golang image as the base image
FROM golang:1.24

# Set the working directory inside the container
WORKDIR /app

# Copy wait-for-it script
COPY wait-for-it.sh .
RUN chmod +x wait-for-it.sh

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o go-notes .

# Expose the port the app runs on
EXPOSE 8080

# Install golang-migrate CLI
RUN wget -O migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz \
    && tar -xzf migrate.tar.gz -C /usr/local/bin \
    && rm migrate.tar.gz \
    && chmod +x /usr/local/bin/migrate

# Command to run the application
CMD ["./go-notes"]