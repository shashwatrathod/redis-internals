package core

import (
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/shashwatrathod/redis-internals/commons"
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
		return commons.WrongNumberOfArgumentsErr(CMD_PING)
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

func handleGet(args []string, c io.ReadWriter) error {
	if len(args) != 1 {
		return commons.WrongNumberOfArgumentsErr(CMD_GET)
	}

	key := args[0]

	val := Get(key)

	// If the Key doesn't exist in the store
	if val == nil {
		c.Write([]byte("$-1\r\n"))
		return nil
	}

	// If the Key exists but the Value is expired.
	if val.expiry != nil && val.expiry.IsExpired() {
		c.Write([]byte("$-1\r\n"))
		return nil
	}

	c.Write(Encode(val.value, false))
	return nil
}

// Parses the arguments to the SET command. Sets the Key with the specified value
// in the datastore.
func handleSet(args []string, c io.ReadWriter) error {

	if len(args) < 2 {
		return commons.WrongNumberOfArgumentsErr(strings.ToLower(CMD_SET))
	}

	// Key and Value are always the 1st and 2nd arguments.
	key, value := args[0], args[1]

	var expiryTime *utils.ExpiryTime = nil

	// Parse rest of the arguments.
	for i := 2; i < len(args); i++ {
		arg := strings.ToLower(args[i])

		switch arg {
		case "ex":
			i++
			if i >= len(args) || expiryTime != nil {
				return commons.SyntaxErr()
			}
			expiryInSeconds, err := strconv.ParseInt(args[i], 10, 64)
			if err != nil {
				return errors.New("ERR value is not an integer or out of range")
			}
			expiryTime = utils.FromExpiryInSeconds(expiryInSeconds)
		case "px":
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

	var val *Value = &Value{
		value:     value,
		valueType: String, // Default to String until other datatypes are implemented.
		expiry:    expiryTime,
	}
	Put(key, val)
	c.Write(Encode("OK", true))
	return nil
}

func handleTtl(args []string, c io.ReadWriter) error {
	return errors.New("unimplemented")
}

// EvalAndRespond processes the specified Redis command and sends the appropriate
// response over the provided network connection.
func EvalAndRespond(cmd *RedisCmd, c io.ReadWriter) error {
	switch cmd.Cmd {
	case CMD_PING:
		return handlePing(cmd.Args, c)
	case CMD_GET:
		return handleGet(cmd.Args, c)
	case CMD_SET:
		return handleSet(cmd.Args, c)
	case CMD_TTL:
		return handleTtl(cmd.Args, c)
	default:
		return commons.UnknownCommandErr(cmd.Cmd, cmd.Args)
	}
}
