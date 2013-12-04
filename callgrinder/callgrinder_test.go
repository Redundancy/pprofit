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

	functionInfo := p.RootFunctions[0]

	if functionInfo.FunctionIdentifier != "test" {
		t.Errorf("Expected root function to have indentifier test: %v", functionInfo.FunctionIdentifier)
	}

	if p.TotalSamples != 100 {
		t.Errorf("Expected total samples to be equal to the first added sample: %v", p.TotalSamples)
	}

	if functionInfo.ExclusiveSamples != 100 {
		t.Errorf("Expected Exclusive Samples to include the function: %v", functionInfo.ExclusiveSamples)
	}

	if functionInfo.InclusiveSamples != 100 {
		t.Errorf("Expected Include Samples to include the function: %v", functionInfo.InclusiveSamples)
	}

}

func TestRootFunctionAddedWhenAddingSampleWithCallers(t *testing.T) {
	t.Skip("TODO")
}

func TestThatRecursiveFunctionsDoNotCreditThemselvesMultipleTimes(t *testing.T) {
	t.Skip("TODO")
}
