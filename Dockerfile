FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY pkg/ cmd/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o raft-demo cmd/main.go 

FROM scratch

WORKDIR /

COPY --from=builder /app/raft-demo /raft-demo

EXPOSE 8080

ENTRYPOINT ["/raft-demo"]
