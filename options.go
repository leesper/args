package args

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
