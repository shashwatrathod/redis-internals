package eval

import (
	"github.com/shashwatrathod/redis-internals/commons"
	"github.com/shashwatrathod/redis-internals/core/resp"
	"github.com/shashwatrathod/redis-internals/core/store"
)

// Evaluates the response to the 'PING' command and responds
// with the results.
//
// Parameters:
//   - args: Arguments passed to the PING command.
func evalPing(args []string, s store.Store) *EvalResult {
	if len(args) >= 2 {
		return &EvalResult{
			Error:    commons.WrongNumberOfArgumentsErr(PING),
			Response: nil,
		}
	}

	var res []byte

	if len(args) == 0 {
		res = resp.Encode("PONG", true)
	} else {
		res = resp.Encode(args[0], false)
	}

	return &EvalResult{
		Response: res,
		Error:    nil,
	}
}
