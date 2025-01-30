package eval

import (
	"github.com/shashwatrathod/redis-internals/commons"
	"github.com/shashwatrathod/redis-internals/core/resp"
	"github.com/shashwatrathod/redis-internals/core/store"
)

// evalGet evaluates the GET command for the Redis server.
// The GET command returns the value of the specified key. If the key does not exist,
// it returns a special nil value.
func evalGet(args []string) *EvalResult {
	if len(args) != 1 {
		return &EvalResult{
			Response: nil,
			Error:    commons.WrongNumberOfArgumentsErr(GET),
		}
	}

	key := args[0]

	val := store.Get(key)

	// If the Key doesn't exist in the store
	if val == nil {
		return &EvalResult{
			Response: []byte("$-1\r\n"),
			Error:    nil,
		}
	}

	// If the Key exists but the Value is expired.
	if val.Expiry != nil && val.Expiry.IsExpired() {
		return &EvalResult{
			Response: []byte("$-1\r\n"),
			Error:    nil,
		}
	}

	return &EvalResult{
		Response: resp.Encode(val.Value, false),
		Error:    nil,
	}
}
