package commandhandler

import (
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/shashwatrathod/redis-internals/commons"
	"github.com/shashwatrathod/redis-internals/core/commands"
	"github.com/shashwatrathod/redis-internals/core/resp"
	"github.com/shashwatrathod/redis-internals/core/store"
	"github.com/shashwatrathod/redis-internals/utils"
)

// Evaluates the response to the 'PING' command and responds
// with the results.
//
// Parameters:
//   - args: Arguments passed to the PING command.
//   - c: the client connection object to respond to.
func handlePing(args []string, c io.ReadWriter) error {
	if len(args) >= 2 {
		return commons.WrongNumberOfArgumentsErr(commands.PING)
	}

	var response []byte

	if len(args) == 0 {
		response = resp.Encode("PONG", true)
	} else {
		response = resp.Encode(args[0], false)
	}

	_, err := c.Write(response)

	return err
}

// handleGet processes the GET command in Redis, which retrieves the value of a specified key.
// If the key does not exist, it returns a nil response. If the key exists but the value is expired,
// it also returns a nil response. Otherwise, it returns the value associated with the key.
//
// Arguments:
//   - args: A slice of strings containing the command arguments. It should contain exactly one element, the key.
//   - c: An io.ReadWriter interface used to write the response.
//
// Returns:
//   - An error if the number of arguments is incorrect, otherwise nil.
func handleGet(args []string, c io.ReadWriter) error {
	if len(args) != 1 {
		return commons.WrongNumberOfArgumentsErr(commands.GET)
	}

	key := args[0]

	val := store.Get(key)

	// If the Key doesn't exist in the store
	if val == nil {
		c.Write([]byte("$-1\r\n"))
		return nil
	}

	// If the Key exists but the Value is expired.
	if val.Expiry != nil && val.Expiry.IsExpired() {
		c.Write([]byte("$-1\r\n"))
		return nil
	}

	c.Write(resp.Encode(val.Value, false))
	return nil
}

// handleSet handles the SET command in Redis, which sets the value of a key.
// If the key already holds a value, it is overwritten, regardless of its type.
//
// Syntax:
// SET key value [EX seconds] [PX milliseconds]
//
// The command supports the following options:
//   - EX seconds: Set the specified expire time, in seconds.
//   - PX milliseconds: Set the specified expire time, in milliseconds.
//
// Example usage:
// SET mykey "Hello"
// SET mykey "Hello" EX 10
// SET mykey "Hello" PX 10000
//
// Parameters:
//   - args: A slice of strings containing the command arguments. The first argument is the key, the second is the value,
//     and optional arguments can specify the expiration time in seconds (EX) or milliseconds (PX).
//   - c: An io.ReadWriter interface for reading from and writing to the client.
//
// Returns:
//   - An error if the command is malformed or if there are issues with parsing the expiration time.
func handleSet(args []string, c io.ReadWriter) error {

	if len(args) < 2 {
		return commons.WrongNumberOfArgumentsErr(strings.ToLower(commands.SET))
	}

	// Key and Value are always the 1st and 2nd arguments.
	key, value := args[0], args[1]

	var expiryTime *utils.ExpiryTime = nil

	// Parse rest of the arguments.
	for i := 2; i < len(args); i++ {
		arg := strings.ToLower(args[i])

		switch arg {
		case commands.EX:
			i++
			if i >= len(args) || expiryTime != nil {
				return commons.SyntaxErr()
			}
			expiryInSeconds, err := strconv.ParseInt(args[i], 10, 64)
			if err != nil {
				return errors.New("ERR value is not an integer or out of range")
			}
			expiryTime = utils.FromExpiryInSeconds(expiryInSeconds)
		case commands.PX:
			i++
			if i >= len(args) || expiryTime != nil {
				return commons.SyntaxErr()
			}
			expiryInMs, err := strconv.ParseInt(args[i], 10, 64)
			if err != nil {
				return errors.New("ERR value is not an integer or out of range")
			}
			expiryTime = utils.FromExpiryInMilliseconds(expiryInMs)
		default:
			return commons.SyntaxErr()
		}
	}

	var val *store.Value = &store.Value{
		Value:     value,
		ValueType: store.String, // Default to String until other datatypes are implemented.
		Expiry:    expiryTime,
	}
	store.Put(key, val)
	c.Write(resp.Encode("OK", true))
	return nil
}

// handleTtl handles the TTL (Time to Live) command for a given key in the Redis store.
// It returns the remaining time to live of a key that has a timeout.
//
// Behavior:
//   - If the key does not exist, it writes -2 to the client.
//   - If the key exists but has no associated expiry, it writes -1 to the client.
//   - If the key exists and has an associated expiry that has not yet passed, it writes the remaining time to live in seconds to the client.
//
// Arguments:
//   - args: A slice of strings where the first element is the key for which TTL is to be checked.
//   - c: An io.ReadWriter interface used for reading from and writing to the client.
//
// Returns:
//   - An error if the number of arguments is incorrect.
func handleTtl(args []string, c io.ReadWriter) error {
	if len(args) != 1 {
		return commons.WrongNumberOfArgumentsErr(commands.GET)
	}

	key := args[0]

	val := store.Get(key)

	// If the Key doesn't exist in the store
	if val == nil {
		c.Write(resp.EncodeWithDatatype(-2, resp.RespInteger))
		return nil
	}

	// The Key exists but there is no expiry associated with it.
	if val.Expiry == nil {
		c.Write(resp.EncodeWithDatatype(-1, resp.RespInteger))
		return nil
	}

	// The expiry has already passed.
	if val.Expiry != nil && val.Expiry.IsExpired() {
		c.Write(resp.EncodeWithDatatype(-2, resp.RespInteger))
		return nil
	}

	c.Write(resp.EncodeWithDatatype(val.Expiry.GetTimeRemainingInSeconds(), resp.RespInteger))
	return nil
}

// EvalAndRespond processes the specified Redis command and sends the appropriate
// response over the provided network connection.
func EvalAndRespond(cmd *commands.RedisCmd, c io.ReadWriter) error {
	switch cmd.Cmd {
	case commands.PING:
		return handlePing(cmd.Args, c)
	case commands.GET:
		return handleGet(cmd.Args, c)
	case commands.SET:
		return handleSet(cmd.Args, c)
	case commands.TTL:
		return handleTtl(cmd.Args, c)
	default:
		return commons.UnknownCommandErr(cmd.Cmd, cmd.Args)
	}
}
