package handlers

func setCommand(args []string) string{
	if len(args) < 2 {
		return "-ERR wrong number of arguments for 'set' command\r\n"
	}

	key := args[0]
	value := args[1]

	mu.Lock()

	store[key] = value

	mu.Unlock()

	return "+OK\r\n"
}