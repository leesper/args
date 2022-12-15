package args

import (
	"reflect"
	"strconv"
)

var supportedFlags = map[string]func(idx int, args []string) interface{}{
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

type BooleanOption bool

func (b BooleanOption) Logging() bool {
	return bool(b)
}

type IntOption int

func (i IntOption) Port() int {
	return int(i)
}

type StringOption string

func (s StringOption) Directory() string {
	return string(s)
}

type MultiOptions struct {
	BooleanOption
	IntOption
	StringOption
}

func Parse(args ...string) interface{} {
	return createOptions(doParse(args))
}

func doParse(args []string) interface{} {
	if len(args) == 0 {
		return false
	}

	values := []interface{}{}

	for flag, fn := range supportedFlags {
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

// factory function for creating options
func createOptions(value interface{}) interface{} {
	// FIXME: replace switch statements with polymorphism
	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Bool:
		return BooleanOption(value.(bool))
	case reflect.Int:
		return IntOption(value.(int))
	case reflect.String:
		return StringOption(value.(string))
	case reflect.Slice:
		singleOptions := map[string]interface{}{
			"bool":   BooleanOption(false),
			"int":    IntOption(0),
			"string": StringOption(""),
		}
		vals := value.([]interface{})
		for _, val := range vals {
			switch val := val.(type) {
			case bool:
				singleOptions["bool"] = BooleanOption(val)
			case int:
				singleOptions["int"] = IntOption(val)
			case string:
				singleOptions["string"] = StringOption(val)
			}
		}
		return MultiOptions{
			BooleanOption: singleOptions["bool"].(BooleanOption),
			IntOption:     singleOptions["int"].(IntOption),
			StringOption:  singleOptions["string"].(StringOption),
		}
	}

	return nil
}
