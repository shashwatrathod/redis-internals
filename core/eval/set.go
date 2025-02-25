package eval

import (
	"errors"
	"strconv"
	"strings"

	"github.com/shashwatrathod/redis-internals/commons"
	"github.com/shashwatrathod/redis-internals/core/resp"
	"github.com/shashwatrathod/redis-internals/core/store"
	"github.com/shashwatrathod/redis-internals/utils"
)

// evalSet processes the SET command with optional arguments to control expiry and insertion.
// Returns an EvalResult with the operation status.
func evalSet(args []string, s store.Store) *EvalResult {
	if len(args) < 2 {
		return &EvalResult{
			Error:    commons.WrongNumberOfArgumentsErr(SET),
			Response: nil,
		}
	}

	// Key and Value are always the 1st and 2nd arguments.
	key, value := args[0], args[1]

	var expiryTime *utils.ExpiryTime = nil

	// Parse rest of the arguments.
	for i := 2; i < len(args); i++ {
		arg := strings.ToLower(args[i])

		switch arg {
		case EX:
			i++
			if i >= len(args) || expiryTime != nil {
				return &EvalResult{
					Error:    commons.SyntaxErr(),
					Response: nil,
				}
			}
			expiryInSeconds, err := strconv.ParseInt(args[i], 10, 64)
			if err != nil {
				return &EvalResult{
					Error:    errors.New("ERR value is not an integer or out of range"),
					Response: nil,
				}
			}
			expiryTime = utils.FromExpiryInSeconds(expiryInSeconds)
		case PX:
			i++
			if i >= len(args) || expiryTime != nil {
				return &EvalResult{
					Error:    commons.SyntaxErr(),
					Response: nil,
				}
			}
			expiryInMs, err := strconv.ParseInt(args[i], 10, 64)
			if err != nil {
				return &EvalResult{
					Error:    errors.New("ERR value is not an integer or out of range"),
					Response: nil,
				}
			}
			expiryTime = utils.FromExpiryInMilliseconds(expiryInMs)
		default:
			return &EvalResult{
				Error:    commons.SyntaxErr(),
				Response: nil,
			}
		}
	}

	s.Put(key, value, expiryTime)

	return &EvalResult{
		Response: resp.Encode("OK", true),
		Error:    nil,
	}
}
