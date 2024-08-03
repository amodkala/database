package main

import (
    "log"
    "net/http"
    "os"

    "github.com/amodkala/raft/pkg/handler"
    "github.com/amodkala/raft/pkg/raft"
)

const (
    raftAddressVar = "RAFT_ADDRESS"
    raftAddressDefault = "localhost:8081"

    serverAddressVar = "SERVER_ADDRESS"
    serverAddressDefault = "localhost:8080"
)

func main() {

    raftAddress, ok := os.LookupEnv(raftAddressVar)
    if !ok {
        raftAddress = raftAddressDefault
    }

    serverAddress, ok := os.LookupEnv(serverAddressVar)
    if !ok {
        serverAddress = serverAddressDefault
    }

    cm := raft.New("raft", serverAddress)

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
