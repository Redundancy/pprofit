/*
Package callgrinder takes lists of samples and constructs a callgraph, with exclusive and inclusive function costs
*/
package callgrinder

import (
	"time"
)

type CPUProfile struct {
	TotalSamples    uint64
	ProfileDuration time.Duration

	// functions at the root of callstacks.
	// these are typically the program entrypoint and goroutine entrypoints
	RootFunctions []FunctionSummary
}

type FunctionSummary struct {
	// An identifier for the function (probably including module and potentially class/type name for methods)
	FunctionIdentifier string

	// The number of samples in this function
	ExclusiveSamples uint64
	// The number of samples in this function, and functions called by this function
	InclusiveSamples uint64

	// TODO: figure out what to do with caller / callee (should it have percentages?)

	// FunctionSummaries that were called by this one
	Callees []FunctionSummary
	// FunctionSummaries that called this one
	Callers []FunctionSummary
}
