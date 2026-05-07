package server

import (
	"bufio"
	"net"
	"testing"
	"time"

	"github.com/kuplays/basiccachedb/internal/executor"
	"github.com/kuplays/basiccachedb/internal/storage"
)

func newTestServer() *Server {
	s := storage.NewStorage()
	e := executor.NewExecutor(s)
	return NewServer(":0", e)
}

func TestServerPing(t *testing.T) {
	srv := newTestServer()

	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	go srv.HandleConnection(serverConn)

	_, err := clientConn.Write([]byte("PING\n"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	reader := bufio.NewReader(clientConn)

	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response != "pong\n" {
		t.Fatalf("expected PONG, got %q", response)
	}
}

func TestServerSetAndGet(t *testing.T) {
	srv := newTestServer()

	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	go srv.HandleConnection(serverConn)

	reader := bufio.NewReader(clientConn)

	_, err := clientConn.Write([]byte("SET name Alice\n"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response != "ok\n" {
		t.Fatalf("expected OK, got %q", response)
	}

	_, err = clientConn.Write([]byte("GET name\n"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	response, err = reader.ReadString('\n')
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response != "Alice\n" {
		t.Fatalf("expected Alice, got %q", response)
	}
}

func TestServerUnknownCommand(t *testing.T) {
	srv := newTestServer()

	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	go srv.HandleConnection(serverConn)

	reader := bufio.NewReader(clientConn)

	_, err := clientConn.Write([]byte("UNKNOWN\n"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response != "ERROR unknown command\n" {
		t.Fatalf("expected unknown command error, got %q", response)
	}
}

func TestServerQuit(t *testing.T) {
	srv := newTestServer()

	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	go srv.HandleConnection(serverConn)

	reader := bufio.NewReader(clientConn)

	_, err := clientConn.Write([]byte("QUIT\n"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response != "OK Quiting...\n" {
		t.Fatalf("expected OK bye, got %q", response)
	}

	_ = clientConn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

	_, err = reader.ReadString('\n')
	if err == nil {
		t.Fatal("expected connection to be closed")
	}
}
