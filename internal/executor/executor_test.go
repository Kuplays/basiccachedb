package executor

import (
	"testing"

	"github.com/kuplays/basiccachedb/internal/protocol"
	"github.com/kuplays/basiccachedb/internal/storage"
)

func newTestExecutor() *Executor {
	s := storage.NewStorage()
	return NewExecutor(s)
}

func TestExecutorSetAndGet(t *testing.T) {
	e := newTestExecutor()

	response := e.Execute(protocol.Command{
		Name: "SET",
		Args: []string{"name", "Alice"},
	})

	if response != RespOK {
		t.Fatalf("expected %s, got %s", RespOK, response)
	}

	response = e.Execute(protocol.Command{
		Name: "GET",
		Args: []string{"name"},
	})

	if response != "Alice" {
		t.Fatalf("expected Alice, got %s", response)
	}
}

func TestExecutorGetMissingKey(t *testing.T) {
	e := newTestExecutor()

	response := e.Execute(protocol.Command{
		Name: "GET",
		Args: []string{"missing"},
	})

	if response != RespNil {
		t.Fatalf("expected %s, got %s", RespNil, response)
	}
}

func TestExecutorDeleteExistingKey(t *testing.T) {
	e := newTestExecutor()

	_ = e.Execute(protocol.Command{
		Name: "SET",
		Args: []string{"name", "Alice"},
	})

	response := e.Execute(protocol.Command{
		Name: "DEL",
		Args: []string{"name"},
	})

	if response != "1" {
		t.Fatalf("expected 1, got %s", response)
	}

	response = e.Execute(protocol.Command{
		Name: "GET",
		Args: []string{"name"},
	})

	if response != RespNil {
		t.Fatalf("expected %s after delete, got %s", RespNil, response)
	}
}

func TestExecutorDeleteMissingKey(t *testing.T) {
	e := newTestExecutor()

	response := e.Execute(protocol.Command{
		Name: "DEL",
		Args: []string{"missing"},
	})

	if response != "0" {
		t.Fatalf("expected 0, got %s", response)
	}
}

func TestExecutorExists(t *testing.T) {
	e := newTestExecutor()

	response := e.Execute(protocol.Command{
		Name: "EXISTS",
		Args: []string{"name"},
	})

	if response != "0" {
		t.Fatalf("expected 0, got %s", response)
	}

	_ = e.Execute(protocol.Command{
		Name: "SET",
		Args: []string{"name", "Alice"},
	})

	response = e.Execute(protocol.Command{
		Name: "EXISTS",
		Args: []string{"name"},
	})

	if response != "1" {
		t.Fatalf("expected 1, got %s", response)
	}
}

func TestExecutorPing(t *testing.T) {
	e := newTestExecutor()

	response := e.Execute(protocol.Command{
		Name: "PING",
		Args: []string{},
	})

	if response != RespPong {
		t.Fatalf("expected %s, got %s", RespPong, response)
	}
}

func TestExecutorSize(t *testing.T) {
	e := newTestExecutor()

	response := e.Execute(protocol.Command{
		Name: "SIZE",
		Args: []string{},
	})

	if response != "0" {
		t.Fatalf("expected 0, got %s", response)
	}

	_ = e.Execute(protocol.Command{
		Name: "SET",
		Args: []string{"a", "1"},
	})

	_ = e.Execute(protocol.Command{
		Name: "SET",
		Args: []string{"b", "2"},
	})

	response = e.Execute(protocol.Command{
		Name: "SIZE",
		Args: []string{},
	})

	if response != "2" {
		t.Fatalf("expected 2, got %s", response)
	}
}

func TestExecutorSetWrongArgCount(t *testing.T) {
	e := newTestExecutor()

	response := e.Execute(protocol.Command{
		Name: "SET",
		Args: []string{"name"},
	})

	if response != RespWrongArgsCount {
		t.Fatalf("expected %s, got %s", RespWrongArgsCount, response)
	}
}

func TestExecutorGetWrongArgCount(t *testing.T) {
	e := newTestExecutor()

	response := e.Execute(protocol.Command{
		Name: "GET",
		Args: []string{},
	})

	if response != RespWrongArgsCount {
		t.Fatalf("expected %s, got %s", RespWrongArgsCount, response)
	}
}

func TestExecutorDelWrongArgCount(t *testing.T) {
	e := newTestExecutor()

	response := e.Execute(protocol.Command{
		Name: "DEL",
		Args: []string{},
	})

	if response != RespWrongArgsCount {
		t.Fatalf("expected %s, got %s", RespWrongArgsCount, response)
	}
}

func TestExecutorExistsWrongArgCount(t *testing.T) {
	e := newTestExecutor()

	response := e.Execute(protocol.Command{
		Name: "EXISTS",
		Args: []string{},
	})

	if response != RespWrongArgsCount {
		t.Fatalf("expected %s, got %s", RespWrongArgsCount, response)
	}
}

func TestExecutorPingWrongArgCount(t *testing.T) {
	e := newTestExecutor()

	response := e.Execute(protocol.Command{
		Name: "PING",
		Args: []string{"extra"},
	})

	if response != RespWrongArgsCount {
		t.Fatalf("expected %s, got %s", RespWrongArgsCount, response)
	}
}

func TestExecutorSizeWrongArgCount(t *testing.T) {
	e := newTestExecutor()

	response := e.Execute(protocol.Command{
		Name: "SIZE",
		Args: []string{"extra"},
	})

	if response != RespWrongArgsCount {
		t.Fatalf("expected %s, got %s", RespWrongArgsCount, response)
	}
}

func TestExecutorUnknownCommand(t *testing.T) {
	e := newTestExecutor()

	response := e.Execute(protocol.Command{
		Name: "UNKNOWN",
		Args: []string{},
	})

	if response != RespUnknownCommand {
		t.Fatalf("expected %s, got %s", RespUnknownCommand, response)
	}
}

func TestExecutorEmptyKey(t *testing.T) {
	e := newTestExecutor()

	response := e.Execute(protocol.Command{
		Name: "SET",
		Args: []string{"", "value"},
	})

	if response != RespInvalidArg {
		t.Fatalf("expected %s, got %s", RespInvalidArg, response)
	}
}
