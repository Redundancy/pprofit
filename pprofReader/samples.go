package pprofReader

import (
	"io"
)

type ProfileSamples []Sample

type Sample struct {
	// number of times the sample was collected
	SampleCount uint64

	// the function pointer stack
	CallStack []uint64
}

// returns unique function pointers (not nessecarily unique functions)
func (p ProfileSamples) UniqueFunctions() []uint64 {
	seenFunctions := make(map[uint64]bool)
	functions := []uint64{}

	for _, sample := range p {
		for _, functionPointer := range sample.CallStack {
			if !seenFunctions[functionPointer] {
				functions = append(functions, functionPointer)
				seenFunctions[functionPointer] = true
			}
		}
	}

	return functions
}

// See documentation on the Binary Trailer of the format
func (s Sample) IsEnd() bool {
	return s.SampleCount == 0 && len(s.CallStack) == 1 && s.CallStack[0] == 0
}

func read32BitSample(r io.Reader) (Sample, error) {
	return readSample(
		r,
		read32BitsAsUInt64)
}

func read64BitSample(r io.Reader) (Sample, error) {
	return readSample(
		r,
		read64BitsAsUInt64)
}

func readSample(r io.Reader, reader func(io.Reader) (uint64, error)) (Sample, error) {
	s := Sample{}
	i, err := reader(r)
	s.SampleCount = i

	i, err = reader(r)

	s.CallStack = make([]uint64, 0, i)

	for x := 0; x < cap(s.CallStack); x++ {
		i, err = reader(r)
		s.CallStack = append(s.CallStack, i)
	}

	return s, err
}
