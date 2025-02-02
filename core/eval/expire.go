package eval

import (
	"strconv"

	"github.com/shashwatrathod/redis-internals/commons"
	"github.com/shashwatrathod/redis-internals/core/resp"
	"github.com/shashwatrathod/redis-internals/core/store"
	"github.com/shashwatrathod/redis-internals/utils"
)

func ttlNotSetResponse() *EvalResult {
	return &EvalResult{
		Response: resp.Encode(0, false),
		Error:    nil,
	}
}

func ttlSetResponse() *EvalResult {
	return &EvalResult{
		Response: resp.Encode(1, false),
		Error:    nil,
	}
}

// Evaluates the EXPIRE command by setting the TTL to the given value
// on the provided key.
func evalExpire(args []string, s store.Store) *EvalResult {
	if len(args) < 2 {
		return &EvalResult{
			Error:    commons.WrongNumberOfArgumentsErr(EXPIRE),
			Response: nil,
		}
	}

	key := args[0]
	expiryInSeconds, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return &EvalResult{
			Error:    commons.UnknownCommandErr(EXPIRE, args),
			Response: nil,
		}
	}

	val := s.Get(key)

	if val == nil {
		return ttlNotSetResponse()
	}

	currentExpiry := val.Expiry

	if currentExpiry == nil || val.Expiry.IsExpired() {
		return ttlNotSetResponse()
	}

	proposedExpiry := utils.FromExpiryInSeconds(expiryInSeconds)

	val.Expiry = proposedExpiry

	// TODO: Add support for additional arguments

	return ttlSetResponse()
}
