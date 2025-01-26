package core

type RedisCmd struct {
	// Represents a Redis Command such as "PING", "GET", "SET", etc.
	Cmd string
	// Represents the arguments to the given Cmd
	Args []string
}

const (
	CMD_PING = "PING"
	CMD_GET  = "GET"
	CMD_TTL  = "TTL"
	CMD_SET  = "SET"
)
