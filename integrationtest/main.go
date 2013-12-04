/*
integration test is designed to run the sample app and extract information from it.
*/
package main

import (
	"fmt"
	"github.com/Redundancy/pprofit/httpFetcher"
	"github.com/Redundancy/pprofit/pprofReader"
	"log"
	"os/exec"
	"time"
)

func main() {
	var err error
	cmd := exec.Command("sampleapp.exe")

	if err != nil {
		log.Fatalf("Could not find and run sampleapp: %v", err)
	}

	if err = cmd.Start(); err != nil {
		log.Fatalf("failed to start sample app: %v", err)
	}
	defer cmd.Process.Kill()

	// Needed to just allow the sample app to start up
	// Better would be to check the status (http and exit code) until it's ready
	time.Sleep(time.Second)

	log.Println("Starting profile collection (this will take ~30 secs)")
	fetcher := &httpFetcher.HttpFetcher{}
	buffer, err := fetcher.GetProfile()
	log.Println("Done collecting profile")

	if err != nil {
		log.Fatal(err)
	}

	profile, err := pprofReader.ReadProfile(buffer)
	if err != nil {
		log.Fatalf("Error reading the profile: %v", err)
	}

	functions := profile.Samples.UniqueFunctions()
	pointerToNameMap, err := fetcher.GetFunctionNames(functions)

	if err != nil {
		log.Fatalf("Error getting function name symbols for profile: %v", err)
	}

	PrintSamples(profile.Samples, pointerToNameMap)
}

func PrintSamples(s pprofReader.ProfileSamples, functionMap map[uint64]string) {
	for _, sample := range s {
		fmt.Print(sample.SampleCount, " [")

		for _, functionPointer := range sample.CallStack {
			fmt.Print(functionMap[functionPointer], ", ")
		}
		fmt.Print("]\n")
	}
}
