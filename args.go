package args

import (
	"reflect"
	"strconv"
)

// Parse parses command line flags into options
func Parse(args ...string) interface{} {
	return createOptions(doParse(args))
}

var singleValuedParsers map[string]func(idx int, args []string) interface{}

func init() {
	singleValuedParsers = map[string]func(idx int, args []string) interface{}{
		"-l": func(idx int, args []string) interface{} {
			return true
		},
		"-p": func(idx int, args []string) interface{} {
			val, _ := strconv.Atoi(args[idx+1])
			return val
		},
		"-d": func(idx int, args []string) interface{} {
			return args[idx+1]
		},
	}
}

// factory function for creating options
func createOptions(value interface{}) interface{} {
	return factories[reflect.ValueOf(value).Kind()].create(value)
}

func doParse(args []string) interface{} {
	if len(args) == 0 {
		return false
	}

	values := []interface{}{}

	for flag, fn := range singleValuedParsers {
		val := parse(flag, args, fn)
		if val != nil {
			values = append(values, val)
		}
	}

	if len(values) == 1 {
		return values[0]
	}

	return values
}

func parse(flag string, args []string, valFn func(idx int, args []string) interface{}) interface{} {
	if idx := contains(args, flag); idx != -1 {
		return valFn(idx, args)
	}
	return nil
}

func contains(args []string, flag string) int {
	for idx, arg := range args {
		if arg == flag {
			return idx
		}
	}
	return -1
}
