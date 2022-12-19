package args

import (
	"reflect"
	"strconv"
)

var singleValuedParsers = map[string]func(idx int, args []string) interface{}{
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

type optionFactory interface {
	create(value interface{}) interface{}
}

type factoryFn func(interface{}) interface{}

func (cf factoryFn) create(value interface{}) interface{} {
	return cf(value)
}

var factories = map[reflect.Kind]optionFactory{
	reflect.Bool:   factoryFn(booleanOptionFactory),
	reflect.Int:    factoryFn(intOptionFactory),
	reflect.String: factoryFn(stringOptionFactory),
	reflect.Slice:  factoryFn(multiOptionsFactory),
}

// FIXME: duplicate code
func booleanOptionFactory(value interface{}) interface{} {
	return BooleanOption(value.(bool))
}

func intOptionFactory(value interface{}) interface{} {
	return IntOption(value.(int))
}

func stringOptionFactory(value interface{}) interface{} {
	return StringOption(value.(string))
}

func multiOptionsFactory(value interface{}) interface{} {
	singleValuedOptions := map[string]interface{}{
		"bool":   BooleanOption(false),
		"int":    IntOption(0),
		"string": StringOption(""),
	}
	vals := value.([]interface{})
	for _, val := range vals {
		switch val := val.(type) {
		case bool:
			singleValuedOptions["bool"] = BooleanOption(val)
		case int:
			singleValuedOptions["int"] = IntOption(val)
		case string:
			singleValuedOptions["string"] = StringOption(val)
		}
	}
	return MultiOptions{
		BooleanOption: singleValuedOptions["bool"].(BooleanOption),
		IntOption:     singleValuedOptions["int"].(IntOption),
		StringOption:  singleValuedOptions["string"].(StringOption),
	}
}

// factory function for creating options
func createOptions(value interface{}) interface{} {
	return factories[reflect.ValueOf(value).Kind()].create(value)
}
