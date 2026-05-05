package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/kuplays/basiccachedb/internal/executor"
	"github.com/kuplays/basiccachedb/internal/protocol"
	"github.com/kuplays/basiccachedb/internal/storage"
)

func main() {
	storage := storage.NewStorage()
	executor := executor.NewExecutor(storage)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Interactive BasicCacheDB example")
	fmt.Println("Type commands: SET key value, GET key, DEL key, EXISTS key, SIZE, PING")
	fmt.Println("Type QUIT to quit")

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()

		if input == "QUIT" {
			break
		}

		cmd, err := protocol.ParseCommand(input)
		if err != nil {
			fmt.Println("ERR", err.Error())
			continue
		}

		response := executor.Execute(cmd)
		fmt.Println(response)
	}
}
