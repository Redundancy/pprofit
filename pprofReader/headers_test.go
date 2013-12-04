package pprofReader

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
	"time"
)

func MakeBufferFrom32BitInts(i ...int32) io.Reader {
	// allocate capacity upfront
	input := bytes.NewBuffer(make([]byte, 0, 40))

	for _, d := range i {
		err := binary.Write(input, binary.LittleEndian, d)
		if err != nil {
			panic(err)
		}
	}

	return input
}

func Test64BitLEHeaderReading(t *testing.T) {

	input := MakeBufferFrom32BitInts(0x00000, 0x00000, 0x00003, 0x00000, 0x00000, 0x00000, 0x02710, 0x00000, 0x00000, 0x00000)
	header, err := profileHeadersFromReader(input)

	if err != nil {
		t.Fatal(err)
	}

	if !header.is64Bit {
		t.Fatal("Header did not detect 64bit header data correctly")
	}

	if header.samplingPeriod != time.Duration(10000)*time.Microsecond {
		t.Errorf("Header did not detect the sampling period of 10000 correctly: %v", header.samplingPeriod)
	}

	if header.version != 0 {
		t.Errorf("Header did not detect the version number correctly: %v", header.version)
	}
}

func Test32BitHeaderReading(t *testing.T) {
	input := MakeBufferFrom32BitInts(0x00000, 0x00003, 0x00000, 0x02710, 0x00000)
	header, err := profileHeadersFromReader(input)

	if err != nil {
		t.Fatal(err)
	}

	if header.is64Bit {
		t.Fatal("Incorrectly detected a 64bit header for a 32bit input")
	}
}
