package main

import (
	"flag"
	"log"

	"github.com/kuplays/basiccachedb/internal/executor"
	"github.com/kuplays/basiccachedb/internal/server"
	"github.com/kuplays/basiccachedb/internal/storage"
)

func main() {
	addr := flag.String("addr", server.DefaultAddress, "TCP addr to listen to")
	flag.Parse()

	storage := storage.NewStorage()
	executor := executor.NewExecutor(storage)

	server := server.NewServer(*addr, executor)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
