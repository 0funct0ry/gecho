# gecho

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`gecho` is a **multi-threaded TCP echo server written in Go**, designed to test connectivity and communication in environments like **Kubernetes clusters**.

It listens for incoming TCP connections and responds with `ECHO [<request-line>]`. If the client sends the command `quit`, the server responds with `Bye` and closes the connection.

---

## 🔧 Features

- Multi-threaded handling of client connections
- Simple echo protocol: responds with `ECHO [<message>]`
- Recognizes `help` command to display available commands and their usage
- Recognizes `quit` command and gracefully closes the client connection
- CLI built with [Cobra](https://github.com/spf13/cobra)
- Customizable bind address and port
- Verbose logging support
- Shell autocompletion generation

---

## 📦 Installation

Clone the repo and build the binary:

```bash
git clone https://github.com/0funct0ry/gecho.git
cd gecho
go build -o gecho
```

---

## 🧰 Command-line Usage

### General Help

```bash
./gecho -h
```

```
An echo server that can be used to test connectivity to a Kubernetes cluster over tcp.

Usage:
  gecho [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  start       Start the echo server

Flags:
  -h, --help   help for gecho

Use "gecho [command] --help" for more information about a command.
```

---

### Start the Echo Server

```bash
./gecho start
```

#### Description

The `start` command launches an echo server that listens for incoming TCP connections  
and echoes back any received data to the client. This server is particularly useful for  
testing network connectivity and validating communication paths in Kubernetes clusters.

#### Flags

| Flag                      | Description                            | Default     |
|---------------------------|----------------------------------------|-------------|
| `-p`, `--port`            | Port number to listen on               | `8080`      |
| `-i`, `--interface`       | The interface to bind to               | `0.0.0.0`   |
| `-v`, `--verbose`         | Enable verbose logging                 | `false`     |
| `-h`, `--help`            | Show help for the `start` command      |             |

#### Example

```bash
./gecho start --interface 127.0.0.1 --port 9000 --verbose
```

---

## 🔄 Sample Client Session

Connect with `nc` or `telnet`:

```bash
nc localhost 8080
```

Example interaction:

```
help
Available commands:
  help       - Display this help message
  quit       - Close the connection
  <message>  - Echo back the message

hello world
ECHO [hello world]

ping
ECHO [ping]

quit
Bye
```

---

## 🧪 Shell Autocompletion

Generate completion script for your shell:

```bash
./gecho completion bash   > /etc/bash_completion.d/gecho
./gecho completion zsh    > ~/.zsh/completions/_gecho
./gecho completion fish   > ~/.config/fish/completions/gecho.fish
./gecho completion powershell > gecho.ps1
```

---

## 📝 License

This project is licensed under the [MIT License](LICENSE).

---

## 📁 Repository

GitHub: [https://github.com/0funct0ry/gecho.git](https://github.com/0funct0ry/gecho.git)

Pull requests and issues are welcome!
