package server

import (
	"bufio"
	"fmt"
	"gecho/commands"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type EchoServer struct {
	port    int
	host    string
	running bool
	ln      net.Listener
	wg      sync.WaitGroup
	exec    *commands.CommandExecutor
}

func NewEchoServer(host string, port int) *EchoServer {
	return &EchoServer{
		port:    port,
		host:    host,
		running: false,
		exec:    commands.NewCommandExecutor(),
	}
}

func (s *EchoServer) Start() error {
	if s.running {
		return fmt.Errorf("server is already running")
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	s.ln = ln
	s.running = true
	log.Printf("Echo server listening on %s:%d", s.host, s.port)

	for s.running {
		conn, err := ln.Accept()
		if err != nil {
			if s.running {
				log.Printf("Failed to accept connection: %v", err)
			}
			continue
		}
		s.wg.Add(1)
		go s.handleConnection(conn)
	}

	return nil
}

func (s *EchoServer) Stop() error {
	if !s.running {
		return nil
	}
	s.running = false
	if err := s.ln.Close(); err != nil {
		return fmt.Errorf("failed to close listener: %v", err)
	}
	s.wg.Wait()
	return nil
}

func (s *EchoServer) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer func(conn net.Conn) { _ = conn.Close() }(conn)

	log.Printf("New connection from %s", conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from connection: %v", err)
			}
			return
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			continue
		}

		if !s.exec.Execute(conn, line) {
			break
		}

	}
}
