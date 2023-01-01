package args

import (
	"reflect"
	"testing"
)

// -l -p 8080 -d /usr/logs
// single option:
func TestShouldSetBooleanOptionToTrueIfFlagPresent(t *testing.T) {
	opt, err := Parse("-l")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	boolOpt := opt.(BooleanOption)
	if !boolOpt.Logging() {
		t.Fatalf("got %t, want %t", boolOpt.Logging(), true)
	}
}

func TestShouldSetBooleanOptionToFalseIfFlagNotPresent(t *testing.T) {
	opt, err := Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	boolOpt := opt.(BooleanOption)
	if boolOpt.Logging() {
		t.Fatalf("got %t, want %t", boolOpt.Logging(), false)
	}
}
func TestShouldParseIntAsOptionValue(t *testing.T) {
	opt, err := Parse("-p", "8080")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	intOpt := opt.(IntOption)
	if intOpt.Port() != 8080 {
		t.Fatalf("got %d, want %d", intOpt.Port(), 8080)
	}
}
func TestShouldParseStrAsOptionValue(t *testing.T) {
	opt, err := Parse("-d", "/usr/logs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	strOpt := opt.(StringOption)
	if strOpt.Directory() != "/usr/logs" {
		t.Fatalf("got %s, want %s", strOpt.Directory(), "/usr/logs")
	}
}

// multiple options:
func TestShouldParseMultiOptions(t *testing.T) {
	opt, err := Parse("-l", "-p", "8080", "-d", "/usr/logs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	multiOpts := opt.(MultiOptions)

	if multiOpts.Logging() != true {
		t.Fatalf("options.Logging() = %t, want %t", multiOpts.Logging(), true)
	}

	if multiOpts.Port() != 8080 {
		t.Fatalf("options.Port() = %d, want %d", multiOpts.Port(), 8080)
	}

	if multiOpts.Directory() != "/usr/logs" {
		t.Fatalf("options.Directory() = %s, want %s", multiOpts.Directory(), "/usr/logs")
	}
}

func TestExample2(t *testing.T) {
	t.SkipNow()
	opt, err := Parse("-g", "this", "is", "a", "list", "-d", "1", "2", "-3", "5")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	listOptions := opt.(ListOptions)

	expectedGroup := []string{"this", "is", "a", "list"}
	if !reflect.DeepEqual(listOptions.Group(), expectedGroup) {
		t.Fatalf("options.Group() = %v, want %v", listOptions.Group(), expectedGroup)
	}

	expectedDecimals := []int{1, 2, -3, 5}
	if !reflect.DeepEqual(listOptions.Decimals(), expectedDecimals) {
		t.Fatalf("options.Decimals() = %v, want %v", listOptions.Decimals(), expectedDecimals)
	}
}

type ListOptions struct{}

func (l ListOptions) Group() []string {
	return nil
}

func (l ListOptions) Decimals() []int {
	return nil
}
