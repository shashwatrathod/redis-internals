package commandhandler

import (
	"io"

	"github.com/shashwatrathod/redis-internals/commons"
	"github.com/shashwatrathod/redis-internals/core/eval"
	"github.com/shashwatrathod/redis-internals/core/store"
)

// EvalAndRespond processes the specified Redis command and sends the appropriate
// response over the provided network connection.
func EvalAndRespond(cmd *eval.RedisCmd, s *store.Store, c io.ReadWriter) error {
	var command *eval.Command = eval.CommandMap[cmd.Cmd]

	if command == nil || command.Eval == nil {
		return commons.UnknownCommandErr(cmd.Cmd, cmd.Args)
	}

	evalResult := command.Eval(cmd.Args, s)

	if evalResult.Error != nil {
		return evalResult.Error
	}

	c.Write(evalResult.Response)
	return nil
}
