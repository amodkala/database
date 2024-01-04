#!/bin/bash

go run cmd/main.go localhost:8080 localhost:8081 localhost:8082 localhost:8083 localhost:8084 &
go run cmd/main.go localhost:8081 localhost:8080 localhost:8082 localhost:8083 localhost:8084 &
go run cmd/main.go localhost:8082 localhost:8081 localhost:8080 localhost:8083 localhost:8084 &
go run cmd/main.go localhost:8083 localhost:8081 localhost:8082 localhost:8080 localhost:8084 &
go run cmd/main.go localhost:8084 localhost:8081 localhost:8082 localhost:8083 localhost:8080
