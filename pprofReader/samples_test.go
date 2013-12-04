package pprofReader

import (
	"testing"
)

// this is currently the only sample in the documentation
// expect this to be fleshed out more by running it against real data
func Test32BitBinaryProfile(t *testing.T) {
	input := MakeBufferFrom32BitInts(5, 3, 0xa0000, 0xc0000, 0xe0000)
	sample, err := read32BitSample(input)

	if err != nil {
		t.Fatal(err)
	}

	if sample.SampleCount != 5 {
		t.Errorf("Expected a sample count of 5, got %v", sample.SampleCount)
	}

	if len(sample.CallStack) != 3 {
		t.Fatalf("Expected a callstack length of 3, got %v", sample.CallStack)
	}

	for i, v := range []uint64{0xa0000, 0xc0000, 0xe0000} {
		if sample.CallStack[i] != v {
			t.Errorf("CallStack[%v] does not match expected value got %v, expected %v ", i, sample.CallStack[i], v)
		}
	}
}
