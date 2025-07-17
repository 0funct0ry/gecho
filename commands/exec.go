package commands

import (
	"fmt"
	"net"
	"strings"
)

type Handler func(conn net.Conn, args []string) bool

type StatusProvider func() map[string]interface{}

type CommandExecutor struct {
	cmdMap        map[string]Handler
	statusProvider StatusProvider
}

func NewCommandExecutor(statusProvider StatusProvider) *CommandExecutor {
	executor := &CommandExecutor{
		cmdMap: map[string]Handler{
			"quit": quitHandler,
			"help": helpHandler,
		},
		statusProvider: statusProvider,
	}

	if statusProvider != nil {
		executor.cmdMap["status"] = executor.statusHandler
	}

	return executor
}

func (e CommandExecutor) Execute(conn net.Conn, line string) bool {
	tokens := strings.Split(line, " ")
	cmd := tokens[0]
	args := tokens[1:]
	handler, ok := e.cmdMap[cmd]
	if !ok {
		return echoHandler(conn, line)
	}
	return handler(conn, args)
}

func helpHandler(conn net.Conn, args []string) bool {
	helpMessage := `
Available commands:
  help       - Display this help message
  status     - Show server uptime and connection statistics
  quit       - Close the connection
  <message>  - Echo back the message
`
	_, _ = conn.Write([]byte(helpMessage))
	return true
}

func quitHandler(conn net.Conn, args []string) bool {
	_, _ = conn.Write([]byte("Bye!\n"))
	return false
}

func echoHandler(conn net.Conn, line string) bool {
	_, _ = conn.Write([]byte(fmt.Sprintf("ECHO [%s]\n", line)))
	return true
}

func (e CommandExecutor) statusHandler(conn net.Conn, args []string) bool {
	if e.statusProvider == nil {
		_, _ = conn.Write([]byte("Status information not available\n"))
		return true
	}

	status := e.statusProvider()

	statusMsg := fmt.Sprintf(`
Server Status:
  Uptime: %s
  Total Connections: %d
  Active Connections: %d
  Start Time: %s
  Address: %s
`,
		status["uptime"],
		status["total_connections"],
		status["active_connections"],
		status["start_time"],
		status["address"])

	_, _ = conn.Write([]byte(statusMsg))
	return true
}
