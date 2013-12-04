/*
This is a very simple sample application for generating data to use against the profile tool

To start with the profile tool is intended to consume data from the net/http/pprof module
This application also has an alternative way to run it that can generate a pprof file (which could be useful for tests)
*/
package main

import (
	"log"
	"math"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/pprof"
)

// rather than using time (and having that show up in the profiles)
// I'm using some approximate loop counters
const (
	LONG  = 1000000
	SHORT = 10000
)

func LongRunningFuncion() {
	for i := 0; i < LONG; i++ {
		math.Sinh(float64(i))
	}
}

func ShortRunningFunction() {
	for i := 0; i < SHORT; i++ {
		math.Sqrt(float64(i))
	}
}

func DoWork() {
	LongRunningFuncion()

	for i := 0; i < SHORT; i++ {
		ShortRunningFunction()
	}
}

func CreateProfile() (f *os.File) {
	var err error
	f, err = os.Create("pprofSample.pprof")

	if err != nil {
		print(err)
		return nil
	}

	pprof.StartCPUProfile(f)

	return f
}

func WaitForTeminationSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}

// Just write a profile to disk
func WriteProfileToDisk() {
	if f := CreateProfile(); f == nil {
		return
	} else {
		defer f.Close()
	}
	DoWork()

	pprof.StopCPUProfile()
}

// Set up a pprof http server, then run in a loop until asked to quit
func HttpProfileServer() {
	go func() {
		log.Println("Starting pprof http server on 6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

Loop:
	for {
		DoWork()

		select {
		case <-c:
			break Loop
		default:
		}
	}

}

func main() {
	writePprofFile := len(os.Args) == 2 && os.Args[1] == "-output"

	if writePprofFile {
		WriteProfileToDisk()
	} else {
		HttpProfileServer()
	}
}
