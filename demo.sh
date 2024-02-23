#!/bin/bash

# Function to kill all child processes upon exit
cleanup() {
    echo "Cleaning up..."
    # Kill the entire process group
    kill 0
}

# Trap SIGINT (Ctrl+C) and SIGTERM
trap cleanup SIGINT SIGTERM

# Start your processes in the background
go run cmd/main.go localhost:8080 localhost:8081 localhost:8082 localhost:8083 localhost:8084 &
go run cmd/main.go localhost:8081 localhost:8080 localhost:8082 localhost:8083 localhost:8084 &
go run cmd/main.go localhost:8082 localhost:8081 localhost:8080 localhost:8083 localhost:8084 &
go run cmd/main.go localhost:8083 localhost:8081 localhost:8082 localhost:8080 localhost:8084 &
go run cmd/main.go localhost:8084 localhost:8081 localhost:8082 localhost:8083 localhost:8080 &

# Wait for all background processes to finish
wait

