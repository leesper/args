package args

import (
	"reflect"
	"testing"
)

// -l -p 8080 -d /usr/logs
// single option:
func TestShouldSetBooleanOptionToTrueIfFlagPresent(t *testing.T) {
	option := Parse("-l").(BooleanOption)
	if !option.Logging() {
		t.Fatalf("got %t, want %t", option.Logging(), true)
	}
}

func TestShouldSetBooleanOptionToFalseIfFlagNotPresent(t *testing.T) {
	option := Parse().(BooleanOption)
	if option.Logging() {
		t.Fatalf("got %t, want %t", option.Logging(), false)
	}
}
func TestShouldParseIntAsOptionValue(t *testing.T) {
	option := Parse("-p", "8080").(IntOption)
	if option.Port() != 8080 {
		t.Fatalf("got %d, want %d", option.Port(), 8080)
	}
}
func TestShouldParseStrAsOptionValue(t *testing.T) {
	option := Parse("-d", "/usr/logs").(StringOption)
	if option.Directory() != "/usr/logs" {
		t.Fatalf("got %s, want %s", option.Directory(), "/usr/logs")
	}
}

// multiple options:
func TestShouldParseMultiOptions(t *testing.T) {
	options := Parse("-l", "-p", "8080", "-d", "/usr/logs").(MultiOptions)

	if options.Logging() != true {
		t.Fatalf("options.Logging() = %t, want %t", options.Logging(), true)
	}

	if options.Port() != 8080 {
		t.Fatalf("options.Port() = %d, want %d", options.Port(), 8080)
	}

	if options.Directory() != "/usr/logs" {
		t.Fatalf("options.Directory() = %s, want %s", options.Directory(), "/usr/logs")
	}
}

// sad path
// 	TODO: - bool -l t / -l t f
//  TODO: - int -p/-p 8080 8081
//  TODO: - string -d / -d /usr/logs /usr/vars
// default values
//  TODO: - bool: false
//  TODO: - integer: 0
//  TODO: - string: ""

func TestExample2(t *testing.T) {
	t.SkipNow()
	listOptions := Parse("-g", "this", "is", "a", "list", "-d", "1", "2", "-3", "5").(ListOptions)

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
