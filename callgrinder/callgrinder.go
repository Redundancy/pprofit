/*
Package callgrinder takes lists of samples and constructs a callgraph, with exclusive and inclusive function costs

Currently WIP
*/
package callgrinder

import (
	"fmt"
	"strings"
	"time"
)

type CPUProfile struct {
	TotalSamples    uint64
	ProfileDuration time.Duration

	// functions at the root of callstacks.
	// these are typically the program entrypoint and goroutine entrypoints
	RootFunctions []*FunctionSummary

	// flattened functions by identifier
	AllFunctions map[string]*FunctionSummary
}

type FunctionSummary struct {
	// An identifier for the function (probably including module and potentially class/type name for methods)
	FunctionIdentifier string

	// The number of samples in this function
	ExclusiveSamples uint64
	// The number of samples in this function, and functions called by this function
	InclusiveSamples uint64

	// TODO: figure out what to do with caller / callee (should it have percentages depending on who called most?)

	// FunctionSummaries that were called by this one
	Callees []*FunctionSummary
	// FunctionSummaries that called this one
	Callers []*FunctionSummary
}

// Add sample information for a callstack
// functions are listed from deepest to shallowest by a string indentifier
func (p *CPUProfile) AddSample(sampleCount uint64, functionIDs ...string) {
	if p.AllFunctions == nil {
		p.AllFunctions = make(map[string]*FunctionSummary)
	}

	if len(functionIDs) == 0 {
		return
	}

	// on recursive functions, don't add them to the inclusive samples more than once
	inclusivelyCreditedFunctions := make(map[string]bool)

	lastFunctionIndex := len(functionIDs) - 1
	var lastFunction *FunctionSummary = nil
	p.TotalSamples += sampleCount

	for i, functionID := range functionIDs {
		functionInfo, existsAlready := p.AllFunctions[functionID]

		if !existsAlready {
			functionInfo = &FunctionSummary{
				FunctionIdentifier: functionID,
			}

			p.AllFunctions[functionID] = functionInfo
		}

		// can be both first and last
		if i == 0 {
			functionInfo.ExclusiveSamples += sampleCount
		}
		if i == lastFunctionIndex && !existsIn(functionInfo, p.RootFunctions) {

			p.RootFunctions = append(p.RootFunctions, functionInfo)

		}

		if lastFunction != nil {
			if !existsIn(functionInfo, lastFunction.Callers) {
				lastFunction.Callers = append(lastFunction.Callers, functionInfo)
			}
			if !existsIn(lastFunction, functionInfo.Callees) {
				functionInfo.Callees = append(functionInfo.Callees, lastFunction)
			}
		}

		if !inclusivelyCreditedFunctions[functionID] {
			functionInfo.InclusiveSamples += sampleCount
			inclusivelyCreditedFunctions[functionID] = true
		}

		lastFunction = functionInfo
	}
}

func existsIn(f *FunctionSummary, l []*FunctionSummary) bool {
	for _, item := range l {
		if item == f {
			return true
		}
	}
	return false
}

// Print a simple summary of a CPU callgraph
// output features inclusive samples and exclusive samples in ()
func (p *CPUProfile) Print() {
	fmt.Printf("Profile with %v samples\n", p.TotalSamples)

	// depth first print
	for _, root := range p.RootFunctions {
		p.printFunction(
			0,
			20,
			root,
			make([]*FunctionSummary, 0, 20),
		)
	}
}

// Could be optimized a bit by not copying the seenFunctionStack each time
func (p *CPUProfile) printFunction(
	depth, limit int,
	f *FunctionSummary,
	seenFunctionStack []*FunctionSummary,
) {

	seenFunctionStack = append(seenFunctionStack, f)

	fmt.Printf(
		"%v %v %.2f%% (%.2f%%)\n",
		strings.Repeat("  ", depth),
		f.FunctionIdentifier,
		100.0*float64(f.InclusiveSamples)/float64(p.TotalSamples),
		100.0*float64(f.ExclusiveSamples)/float64(p.TotalSamples),
	)

	if depth >= limit {
		return
	}

ForEachCallee:
	for _, c := range f.Callees {
		if existsIn(c, seenFunctionStack) {
			fmt.Printf(
				"%v %v %.2f%% (%.2f%%) [R]\n",
				strings.Repeat("  ", depth+1),
				c.FunctionIdentifier,
				100.0*float64(c.InclusiveSamples)/float64(p.TotalSamples),
				100.0*float64(c.ExclusiveSamples)/float64(p.TotalSamples),
			)
			continue ForEachCallee
		}

		p.printFunction(
			depth+1,
			limit,
			c,
			seenFunctionStack, // NB: copy, so no need to remove items
		)
	}
}
