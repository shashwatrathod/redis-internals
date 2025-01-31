package eval

import (
	"github.com/shashwatrathod/redis-internals/commons"
	"github.com/shashwatrathod/redis-internals/core/resp"
	"github.com/shashwatrathod/redis-internals/core/store"
)

// evalDel processes the DEL command and deletes the keys passed in the arguments from the store.
// Returns the number of keys deleted in the result.
func evalDel(args []string) *EvalResult {
	if len(args) == 0 {
		return &EvalResult{
			Error:    commons.WrongNumberOfArgumentsErr(DEL),
			Response: nil,
		}
	}

	nDeleted := 0

	for _, key := range args {
		if isDeleted := store.Delete(key); isDeleted {
			nDeleted++
		}
	}

	return &EvalResult{
		Response: resp.Encode(nDeleted, false),
		Error:    nil,
	}
}
