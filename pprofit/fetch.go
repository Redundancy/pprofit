package main

import (
	"github.com/Redundancy/pprofit/callgrinder"
	"github.com/Redundancy/pprofit/httpFetcher"
	"github.com/Redundancy/pprofit/pprofReader"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"time"
)

func init() {
	app.Commands = append(
		app.Commands,
		cli.Command{
			Name:      "fetch",
			ShortName: "f",
			Usage:     "get profile information over http",
			Action:    Fetch,
			Flags: []cli.Flag{
				cli.StringFlag{"server, s", "localhost", "the server to go to"},
				cli.IntFlag{"port, p", 6060, "The port to use"},
				cli.IntFlag{"duration, d", 30, "Duration of the profile in seconds"},
			},
		},
	)
}

func Fetch(c *cli.Context) {
	fetcher := &httpFetcher.HttpFetcher{
		Server:          c.String("server"),
		Port:            uint(c.Int("port")),
		ProfileDuration: time.Duration(c.Int("duration")) * time.Second,
	}

	buffer, err := fetcher.GetProfile()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	profile, err := pprofReader.ReadProfile(buffer)
	if err != nil {
		log.Fatalf("Error reading the profile: %v", err)
		os.Exit(1)
	}

	functions := profile.Samples.UniqueFunctions()
	pointerToNameMap, err := fetcher.GetFunctionNames(functions)

	if err != nil {
		log.Fatalf("Error getting function name symbols for profile: %v", err)
		os.Exit(1)
	}

	crunchedProfile := &callgrinder.CPUProfile{}
	for _, sample := range profile.Samples {
		crunchedProfile.AddSample(
			sample.SampleCount,
			SampleFunctionsToIdentifiers(sample, pointerToNameMap)...,
		)
	}

	crunchedProfile.Print()
}

func SampleFunctionsToIdentifiers(s pprofReader.Sample, functionMap map[uint64]string) []string {
	result := make([]string, len(s.CallStack))

	for i, pointer := range s.CallStack {
		result[i] = functionMap[pointer]
	}

	return result
}
