package pprofReader

import (
	"fmt"
	"os"
	"testing"
	"time"
)

const standardSamplingPeriod = time.Duration(10000) * time.Microsecond

func TestReadingSampleProfile(t *testing.T) {
	f, err := os.Open("testdata/pprofSample.pprof")

	if err != nil {
		t.Fatal(err)
	}

	profile, err := ReadProfile(f)

	if err != nil {
		t.Fatal(err)
	}

	if profile == nil {
		t.Error("Profile was nil")
	}

	if !profile.headers.is64Bit {
		t.Error("Expected a 64bit profile")
	}

	if profile.headers.samplingPeriod != standardSamplingPeriod {
		t.Errorf("Unexpected sampling period: %v (expected: %v)", profile.headers.samplingPeriod, standardSamplingPeriod)
	}

	if len(profile.Samples) == 0 {
		t.Fatal("No samples?")
	}

	/*
		for _, s := range profile.Samples {
			fmt.Println(s)
		}
	*/
}
