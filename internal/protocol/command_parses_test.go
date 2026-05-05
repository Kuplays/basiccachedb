package protocol

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseCommandSet(t *testing.T) {
	cmd, err := ParseCommand("SET name Alice")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cmd.Name != "SET" {
		t.Fatalf("expected command name SET, got %s", cmd.Name)
	}

	expectedArgs := []string{"name", "Alice"}
	if !reflect.DeepEqual(cmd.Args, expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args)
	}
}

func TestParseCommandGet(t *testing.T) {
	cmd, err := ParseCommand("GET name")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cmd.Name != "GET" {
		t.Fatalf("expected command name GET, got %s", cmd.Name)
	}

	expectedArgs := []string{"name"}
	if !reflect.DeepEqual(cmd.Args, expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args)
	}
}

func TestParseCommandLowercaseName(t *testing.T) {
	cmd, err := ParseCommand("get name")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cmd.Name != "GET" {
		t.Fatalf("expected command name GET, got %s", cmd.Name)
	}
}

func TestParseCommandWithExtraSpaces(t *testing.T) {
	cmd, err := ParseCommand("   SET    name     Alice   ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cmd.Name != "SET" {
		t.Fatalf("expected command name SET, got %s", cmd.Name)
	}

	expectedArgs := []string{"name", "Alice"}
	if !reflect.DeepEqual(cmd.Args, expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, cmd.Args)
	}
}

func TestParseCommandPing(t *testing.T) {
	cmd, err := ParseCommand("PING")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cmd.Name != "PING" {
		t.Fatalf("expected command name PING, got %s", cmd.Name)
	}

	if len(cmd.Args) != 0 {
		t.Fatalf("expected no args, got %v", cmd.Args)
	}
}

func TestParseEmptyCommand(t *testing.T) {
	_, err := ParseCommand("")

	if !errors.Is(err, ErrorEmptyCommand) {
		t.Fatalf("expected ErrEmptyCommand, got %v", err)
	}
}

func TestParseWhitespaceOnlyCommand(t *testing.T) {
	_, err := ParseCommand("      ")

	if !errors.Is(err, ErrorEmptyCommand) {
		t.Fatalf("expected ErrEmptyCommand, got %v", err)
	}
}
