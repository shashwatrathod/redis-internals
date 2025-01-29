package eval

type RedisCmd struct {
	// Represents a Redis Command such as "PING", "GET", "SET", etc.
	Cmd string
	// Represents the arguments to the given Cmd
	Args []string
}

type EvalResponse struct {
	Response []byte
	Error    error
}

type Command struct {
	Name string
	Eval func(args []string) EvalResponse
}

// supported commands
const (
	PING = "PING"
	GET  = "GET"
	TTL  = "TTL"
	SET  = "SET"
)

// supported command arguments
const (
	EX = "ex"
	PX = "px"
)
