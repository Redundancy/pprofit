package callgrinder

import (
	"testing"
)

type ExpectedNameInclusiveExclusive struct {
	name      string
	inclusive uint64
	exclusive uint64
}

func CompareCallstackAgainstExpected(
	functionInfo *FunctionSummary,
	expectedSamples []ExpectedNameInclusiveExclusive,
	t *testing.T,
) {
	for i, expected := range expectedSamples {
		if functionInfo == nil {
			t.Fatal("Unexpected nil function caller for callstack ", i)
		}

		if functionInfo.FunctionIdentifier != expected.name {
			t.Errorf(
				"Function does not match expected name: %v vs expected %v",
				functionInfo.FunctionIdentifier,
				expected.name,
			)
		}
		if functionInfo.ExclusiveSamples != expected.exclusive {
			t.Errorf(
				"Function does not match expected exclusive samples: %v vs expected %v",
				functionInfo.ExclusiveSamples,
				expected.exclusive,
			)
		}
		if functionInfo.InclusiveSamples != expected.inclusive {
			t.Errorf(
				"Function does not match expected inclusive samples: %v vs expected %v",
				functionInfo.InclusiveSamples,
				expected.inclusive,
			)
		}

		if len(functionInfo.Callees) == 1 {
			functionInfo = functionInfo.Callees[0]
		} else if len(functionInfo.Callees) > 1 {
			if i < len(expectedSamples)-1 {
				// Make an attempt at finding the callee that we were supposed to
				found := false
				for _, callee := range functionInfo.Callees {
					if callee.FunctionIdentifier == expectedSamples[i+1].name {
						found = true
						functionInfo = callee
					}
				}
				if !found {
					t.Fatalf(
						"Function was unexpectedly missing a reference to a callee: %v",
						functionInfo.FunctionIdentifier,
					)
				}
			}

		} else {
			// everything except for the last item should call something else
			if i < len(expectedSamples)-1 {
				t.Fatalf(
					"Function was unexpectedly missing a reference to a callee: %v",
					functionInfo.FunctionIdentifier,
				)
			}
		}
	}
}

func TestRootFunctionAddedWhenAddingFirstSampleWithNoCallers(t *testing.T) {
	p := CPUProfile{}

	p.AddSample(100, "test")

	if len(p.RootFunctions) != 1 {
		t.Fatalf("Expected root functions to contain test function: %v", p.RootFunctions)
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
	p := CPUProfile{}
	const samples = 50

	p.AddSample(samples, "test", "testcaller", "testcallercaller")

	if len(p.RootFunctions) != 1 {
		t.Fatalf("Expected root functions to contain one function")
	}

	rootFunctionInfo := p.RootFunctions[0]

	if rootFunctionInfo.FunctionIdentifier != "testcallercaller" {
		t.Errorf("Expected root function to have indentifier testcallercaller: %v", rootFunctionInfo.FunctionIdentifier)
	}

	if p.TotalSamples != samples {
		t.Errorf("Expected total samples to be equal to the first added sample: %v", p.TotalSamples)
	}

	if rootFunctionInfo.ExclusiveSamples != 0 {
		t.Errorf("Root function should not have had any Exclusive samples: %v", rootFunctionInfo.ExclusiveSamples)
	}

	if rootFunctionInfo.InclusiveSamples != samples {
		t.Errorf("Expected Include Samples to include the function: %v", rootFunctionInfo.InclusiveSamples)
	}
}

func TestCallersAreCreated(t *testing.T) {
	p := CPUProfile{}
	const samples = 50

	p.AddSample(samples, "test", "testcaller", "testcallercaller")

	expectedSamples := []ExpectedNameInclusiveExclusive{
		{
			"testcallercaller",
			samples,
			0,
		},
		{
			"testcaller",
			samples,
			0,
		},
		{
			"test",
			samples,
			samples,
		},
	}

	functionInfo := p.RootFunctions[0]
	CompareCallstackAgainstExpected(functionInfo, expectedSamples, t)
}

func TestThatRecursiveFunctionsDoNotCreditThemselvesMultipleTimes(t *testing.T) {
	p := CPUProfile{}
	const samples = 50

	p.AddSample(samples, "test", "testCallerRecursive", "testCallerRecursive")

	expectedSamples := []ExpectedNameInclusiveExclusive{
		{
			"testCallerRecursive",
			samples,
			0,
		},
		{
			"testCallerRecursive",
			samples,
			0,
		},
		{
			"test",
			samples,
			samples,
		},
	}

	functionInfo := p.RootFunctions[0]
	CompareCallstackAgainstExpected(functionInfo, expectedSamples, t)
}

func TestThatCalleesAreNotAddedMultipleTimes(t *testing.T) {
}
