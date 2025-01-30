package eval

type RedisCmd struct {
	// Represents a Redis Command such as "PING", "GET", "SET", etc.
	Cmd string
	// Represents the arguments to the given Cmd
	Args []string
}

// Represents the results of evalutation of a Command.
type EvalResult struct {
	// The resp-encoded byte slice to be returned to the caller.
	Response []byte
	// Any known error that was encountered during the execution.
	Error error
	// TODO: Find a way to send regular responses but also indicate some meta information - like RESP type/etc, so that RESP can be decoupled from core.
}

// Represents a single Redis Command. Knows how to execute the command.
type Command struct {
	// Name of the command (eg. PING, GET, SET, etc.)
	Name string

	// Evaluates the command by executing the core logic and
	// returns the results from the execution.
	Eval func(args []string) *EvalResult
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

// This map will hold data about all the supported Commands.
// Warning : This should be treated as an immutable entity!
var CommandMap = map[string]*Command{}

func init() {
	CommandMap[PING] = &Command{
		Name: PING,
		Eval: evalPing,
	}

	CommandMap[GET] = &Command{
		Name: GET,
		Eval: evalGet,
	}

	CommandMap[SET] = &Command{
		Name: SET,
		Eval: evalSet,
	}

	CommandMap[TTL] = &Command{
		Name: TTL,
		Eval: evalTtl,
	}
}
