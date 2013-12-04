package callgrinder

import (
	"testing"
)

func TestRootFunctionAddedWhenAddingFirstSampleWithNoCallers(t *testing.T) {
	p := CPUProfile{}

	p.AddSample("test", 100)

	if len(p.RootFunctions) != 1 {
		t.Fatalf("Expected root functions to contain test function")
	}

	if p.RootFunctions[0].FunctionIdentifier != "test" {
		t.Errorf("Expected root function to have indentifier test: %v", p.RootFunctions[0].FunctionIdentifier)
	}
}
