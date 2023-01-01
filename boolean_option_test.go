package args

import (
	"testing"
)

// sad path
//   - bool -l t / -l t f
func TestShouldNotAcceptExtraArgumentForBooleanOption(t *testing.T) {
	_, err := doParse([]string{"-l", "t"})
	if err == nil || err != ErrTooManyArgs {
		t.Fatalf("got nil, want %v", ErrTooManyArgs)
	}

}

func TestShouldNotAcceptExtraArgumentsForBooleanOption(t *testing.T) {
	_, err := doParse([]string{"-l", "t", "f"})
	if err == nil || err != ErrTooManyArgs {
		t.Fatalf("got nil, want %v", ErrTooManyArgs)
	}
}

// default values
//   - bool: false
func TestShouldSetDefaultValueToFalseIfOptionNotPresent(t *testing.T) {
	val, err := doParse([]string{})
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	boolVal := val.(bool)
	if boolVal {
		t.Fatalf("got %t, want %t", boolVal, false)
	}
}
