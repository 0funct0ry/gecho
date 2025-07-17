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
	"sync/atomic"
	"time"
)

type EchoServer struct {
	port              int
	host              string
	running           bool
	ln                net.Listener
	wg                sync.WaitGroup
	exec              *commands.CommandExecutor
	startTime         time.Time
	totalConnections  uint64
	activeConnections int32
}

func NewEchoServer(host string, port int) *EchoServer {
	server := &EchoServer{
		port:              port,
		host:              host,
		running:           false,
		totalConnections:  0,
		activeConnections: 0,
	}

	// Create a status provider function that calls GetStatus on this server instance
	statusProvider := func() map[string]interface{} {
		return server.GetStatus()
	}

	server.exec = commands.NewCommandExecutor(statusProvider)

	return server
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
	s.startTime = time.Now()
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

func (s *EchoServer) GetStatus() map[string]interface{} {
	uptime := time.Since(s.startTime)
	return map[string]interface{}{
		"uptime":            uptime.String(),
		"uptime_seconds":    int(uptime.Seconds()),
		"total_connections": atomic.LoadUint64(&s.totalConnections),
		"active_connections": atomic.LoadInt32(&s.activeConnections),
		"start_time":        s.startTime.Format(time.RFC3339),
		"address":           fmt.Sprintf("%s:%d", s.host, s.port),
	}
}

func (s *EchoServer) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer func(conn net.Conn) { 
		atomic.AddInt32(&s.activeConnections, -1)
		_ = conn.Close() 
	}(conn)

	atomic.AddUint64(&s.totalConnections, 1)
	atomic.AddInt32(&s.activeConnections, 1)
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
