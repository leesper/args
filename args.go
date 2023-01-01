package args

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var (
	supportedFlagsParsers map[string]func(idx int, args []string) (interface{}, error)
	ErrTooManyArgs        = errors.New("too many argruments")
)

// Parse parses command line flags into options
func Parse(args ...string) (interface{}, error) {
	value, err := doParse(args)
	if err != nil {
		return value, err
	}
	return createOptions(value)
}

func init() {
	supportedFlagsParsers = map[string]func(idx int, args []string) (interface{}, error){
		"-l": func(idx int, args []string) (interface{}, error) {
			if idx+1 < len(args) && !strings.HasPrefix(args[idx+1], "-") {
				return nil, ErrTooManyArgs
			}
			return true, nil
		},
		"-p": func(idx int, args []string) (interface{}, error) {
			val, _ := strconv.Atoi(args[idx+1])
			return val, nil
		},
		"-d": func(idx int, args []string) (interface{}, error) {
			return args[idx+1], nil
		},
	}
}

// factory function for creating options
func createOptions(value interface{}) (interface{}, error) {
	return factories[reflect.ValueOf(value).Kind()].create(value)
}

func doParse(args []string) (interface{}, error) {
	if len(args) == 0 {
		return false, nil
	}

	values := []interface{}{}

	for flag, fn := range supportedFlagsParsers {
		val, err := parse(flag, args, fn)
		if err != nil {
			return nil, err
		}
		if val != nil {
			values = append(values, val)
		}
	}

	if len(values) == 1 {
		return values[0], nil
	}

	return values, nil
}

func parse(flag string, args []string, valFn func(idx int, args []string) (interface{}, error)) (interface{}, error) {
	if idx := contains(args, flag); idx != -1 {
		return valFn(idx, args)
	}
	return nil, nil
}

func contains(args []string, flag string) int {
	for idx, arg := range args {
		if arg == flag {
			return idx
		}
	}
	return -1
}
