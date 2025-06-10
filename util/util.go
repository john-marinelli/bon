package util

import (
	"errors"

	"github.com/john-marinelli/bon/types"
)

var ArgNumberErr = errors.New("incorrect number of arguments")
var WrongArgErr = errors.New("incorrect arguments")

func ParseArgs(args []string) (types.BonScreen, error) {
	if len(args) == 1 {
		return types.InputScreen, nil
	}
	if len(args) > 2 || len(args) < 1 {
		return types.InputScreen, ArgNumberErr
	}
	if args[1] != "bon" {
		return types.InputScreen, WrongArgErr
	}

	return types.ViewScreen, nil
}

func Truncate(str string, ml int) string {
	r := []rune(str)
	if len(r) > ml {
		return string(r[:ml-3]) + "..."
	}

	return str
}
