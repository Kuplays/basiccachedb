package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/kuplays/basiccachedb/internal/executor"
	"github.com/kuplays/basiccachedb/internal/protocol"
)

const DefaultAddress = ":8000"

type Server struct {
	address  string
	executor *executor.Executor
}

func NewServer(address string, e *executor.Executor) *Server {
	if address == "" {
		address = DefaultAddress
	}

	return &Server{
		address:  address,
		executor: e,
	}
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Printf("BasicCacheDB server is listening on %s\n", s.address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			_, _ = conn.Write([]byte("ERROR read error\n"))
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.EqualFold(line, "QUIT") {
			_, _ = conn.Write([]byte("OK Quiting...\n"))
			return
		}

		cmd, err := protocol.ParseCommand(line)
		if err != nil {
			_, _ = conn.Write([]byte("ERROR " + err.Error() + "\n"))
			continue
		}

		response := s.executor.Execute(cmd)
		_, _ = conn.Write([]byte(response + "\n"))
	}
}
