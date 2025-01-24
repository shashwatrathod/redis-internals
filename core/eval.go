package core

import (
	"io"

	"github.com/shashwatrathod/redis-internals/rediserrors"
)

// Evaluates the response to the 'PING' command and responds
// with the results.
//
// Parameters:
//   - args: Arguments passed to the PING command.
//   - conn: the client connection object to respond to.
func handlePing(args []string, c io.ReadWriter) error {
	if len(args) >= 2 {
		return rediserrors.WrongNumberOfArguments(CmdPing)
	}

	var response []byte

	if len(args) == 0 {
		response = Encode("PONG", true)
	} else {
		response = Encode(args[0], false)
	}

	_, err := c.Write(response)

	return err
}

// EvalAndRespond processes the specified Redis command and sends the appropriate
// response over the provided network connection.
func EvalAndRespond(cmd *RedisCmd, c io.ReadWriter) error {
	switch cmd.Cmd {
	case CmdPing:
		return handlePing(cmd.Args, c)
	default:
		return rediserrors.UnknownCommand(cmd.Cmd, cmd.Args)
	}
}
