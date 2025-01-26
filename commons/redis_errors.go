package commons

import (
	"errors"
	"fmt"
	"strings"
)

func WrongNumberOfArgumentsErr(cmd string) error {
	return fmt.Errorf("ERR wrong number of arguments for '%s' command", strings.ToLower(cmd))
}

func SyntaxErr() error {
	return errors.New("ERR syntax error")
}

func UnknownCommandErr(cmd string, args []string) error {
	argStr := ""
	if len(args) > 0 {
		argStr = fmt.Sprintf("'%s'", strings.Join(args[:], "' '"))
	}
	return fmt.Errorf("ERR unknown command '%s', with args beginning with: %s", cmd, argStr)
}
