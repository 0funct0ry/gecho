package commands

import (
	"fmt"
	"net"
	"strings"
)

type Handler func(conn net.Conn, args []string) bool

type CommandExecutor struct {
	cmdMap map[string]Handler
}

func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{
		cmdMap: map[string]Handler{
			"quit": quitHandler,
			"help": helpHandler,
		},
	}
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
