package core

type RedisCmd struct {
	// Represents a Redis Command such as "PING", "GET", "SET", etc.
	Cmd string
	// Represents the arguments to the given Cmd
	Args []string
}

const (
	CmdPing = "PING"
)
