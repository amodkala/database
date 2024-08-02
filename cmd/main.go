package main

import (
    "log"
    "net/http"

    "github.com/amodkala/raft/pkg/handler"
    "github.com/amodkala/raft/pkg/raft"
)

const (
    name = "example"
    raftAddress = "localhost:8081"
    serverAddress = "localhost:8080"
)

func main() {

    cm := raft.New(name, serverAddress)

    go func() {
        if err := cm.Start(raftAddress); err != nil {
            log.Fatal(err)
        }
    }()    

    writeHandler := handler.NewWriteHandler(cm)

    mux := http.NewServeMux()
    mux.HandleFunc("/write", writeHandler)

    srv := http.Server{
        Addr: serverAddress,
        Handler: mux,
    }

    log.Fatal(srv.ListenAndServe())
}
