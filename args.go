package args

import "strconv"

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
	values := parseArgs(args)
	return createOptions(values)
}

func parseArgs(args []string) []interface{} {
	values := []interface{}{}
	for idx, arg := range args {
		if arg == "-l" {
			values = append(values, true)
		} else if arg == "-p" {
			value, _ := strconv.Atoi(args[idx+1])
			values = append(values, value)
		} else if arg == "-d" {
			values = append(values, args[idx+1])
		}
	}

	if len(values) == 0 {
		values = append(values, false)
	}

	return values
}

// factory function for creating options
func createOptions(values []interface{}) interface{} {
	optionDict := map[string]interface{}{
		"bool":   BooleanOption(false),
		"int":    IntOption(0),
		"string": StringOption(""),
	}

	typeTag := ""
	for _, value := range values {
		switch value := value.(type) {
		case bool:
			typeTag = "bool"
			optionDict[typeTag] = BooleanOption(value)
		case int:
			typeTag = "int"
			optionDict[typeTag] = IntOption(value)
		case string:
			typeTag = "string"
			optionDict[typeTag] = StringOption(value)
		}
	}

	if len(values) == 0 {
		return BooleanOption(false)
	} else if len(values) == 1 {
		return optionDict[typeTag]
	}

	return MultiOptions{
		BooleanOption: optionDict["bool"].(BooleanOption),
		IntOption:     optionDict["int"].(IntOption),
		StringOption:  optionDict["string"].(StringOption),
	}
}
