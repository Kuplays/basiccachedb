package executor

import (
	"errors"
	"fmt"

	"github.com/kuplays/basiccachedb/internal/protocol"
	"github.com/kuplays/basiccachedb/internal/storage"
)

const (
	RespOK             = "ok"
	RespNil            = "nil"
	RespPong           = "pong"
	RespUnknownCommand = "ERROR unknown command"
	RespWrongArgsCount = "ERROR wrong number of arguments"
	RespInternalError  = "ERROR internal error"
	RespInvalidArg     = "ERROR invalid argument"
)

type Executor struct {
	storage *storage.Storage
}

func NewExecutor(s *storage.Storage) *Executor {
	return &Executor{
		storage: s,
	}
}

func (e *Executor) Execute(cmd protocol.Command) string {
	switch cmd.Name {
	case "SET":
		return e.ExecuteSet(cmd.Args)
	case "GET":
		return e.ExecuteGet(cmd.Args)
	case "DEL":
		return e.ExecuteDel(cmd.Args)
	case "EXISTS":
		return e.ExecuteExists(cmd.Args)
	case "PING":
		return e.ExecutePing(cmd.Args)
	case "SIZE":
		return e.ExecuteLength(cmd.Args)
	default:
		return RespUnknownCommand
	}
}

func (e *Executor) ExecuteSet(args []string) string {
	if len(args) != 2 {
		return RespWrongArgsCount
	}

	key := args[0]
	value := args[1]

	err := e.storage.Set(key, value)

	if err != nil {
		if errors.Is(err, storage.ErrorEmptyKey) {
			return RespInvalidArg
		}
		return RespInternalError
	}

	return RespOK
}

func (e *Executor) ExecuteGet(args []string) string {
	if len(args) != 1 {
		return RespWrongArgsCount
	}

	key := args[0]

	value, err := e.storage.Get(key)

	if err != nil {
		if errors.Is(err, storage.ErrorKeyNotFound) {
			return RespNil
		}

		if errors.Is(err, storage.ErrorEmptyKey) {
			return RespInvalidArg
		}

		return RespInternalError
	}

	return value
}

func (e *Executor) ExecuteDel(args []string) string {
	if len(args) != 1 {
		return RespWrongArgsCount
	}

	key := args[0]

	deleted, err := e.storage.Delete(key)

	if err != nil {
		if errors.Is(err, storage.ErrorEmptyKey) {
			return RespInvalidArg
		}

		return RespInternalError
	}

	if deleted {
		return "1"
	}

	return "0"
}

func (e *Executor) ExecuteExists(args []string) string {
	if len(args) != 1 {
		return RespWrongArgsCount
	}

	key := args[0]

	exists, err := e.storage.Exists(key)

	if err != nil {
		if errors.Is(err, storage.ErrorEmptyKey) {
			return RespInvalidArg
		}

		return RespInternalError
	}

	if exists {
		return "1"
	}

	return "0"
}

func (e *Executor) ExecuteLength(args []string) string {
	if len(args) != 0 {
		return RespWrongArgsCount
	}

	return fmt.Sprintf("%d", e.storage.Length())
}

func (e *Executor) ExecutePing(args []string) string {
	if len(args) != 0 {
		return RespWrongArgsCount
	}

	return RespPong
}
