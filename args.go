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

func Parse(args ...string) interface{} {
	for idx, arg := range args {
		if arg == "-l" {
			return BooleanOption(true)
		} else if arg == "-p" {
			value, _ := strconv.Atoi(args[idx+1])
			return IntOption(value)
		} else if arg == "-d" {
			return StringOption(args[idx+1])
		}
	}

	return BooleanOption(false)
}
