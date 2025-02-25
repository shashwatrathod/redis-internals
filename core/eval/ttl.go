package eval

import (
	"github.com/shashwatrathod/redis-internals/commons"
	"github.com/shashwatrathod/redis-internals/core/resp"
	"github.com/shashwatrathod/redis-internals/core/store"
	"github.com/shashwatrathod/redis-internals/utils"
)

// evaluates the TTL (Time to Live) command for a given key in the Redis store.
// It returns the remaining time to live of a key that has a timeout.
func evalTtl(args []string, s store.Store) *EvalResult {
	if len(args) != 1 {
		return &EvalResult{
			Error:    commons.WrongNumberOfArgumentsErr(TTL),
			Response: nil,
		}
	}

	key := args[0]

	val := s.Get(key)

	// If the Key doesn't exist in the store
	if val == nil {
		return &EvalResult{
			Response: resp.EncodeWithDatatype(-2, resp.RespInteger),
			Error:    nil,
		}
	}

	expTs := s.GetExpiry(key)

	// The Key exists but there is no expiry associated with it.
	if expTs == nil {
		return &EvalResult{
			Response: resp.EncodeWithDatatype(-1, resp.RespInteger),
			Error:    nil,
		}
	}

	expiry := utils.FromExpiryInUnixTime(*expTs)

	// The expiry has already passed.
	if expiry != nil && expiry.IsExpired() {
		return &EvalResult{
			Response: resp.EncodeWithDatatype(-2, resp.RespInteger),
			Error:    nil,
		}
	}

	return &EvalResult{
		Response: resp.EncodeWithDatatype(expiry.GetTimeRemainingInSeconds(), resp.RespInteger),
		Error:    nil,
	}
}
