package handlers

import (
	"strings"
)

type Command func(args []string) string

var registry = map[string]Command{
	"PING":pingHandler,
	"ECHO":echoHandler,
	"SET":setCommand,
	"GET":getCommand,
	"RPUSH":rpushCommand,
	"LRANGE":lrangeCommand,
	"LPUSH": lpushCommand,
	"LLEN": llenCommand,
	"LPOP": lpopCommand,
}

func Execute(command string,args []string) (string) {
	command = strings.ToUpper(command)
	if handler,exists := registry[command]; exists {
		return handler(args)
	}
	return "-ERR unknown command '" + command + "'\r\n"
}

