package pprofReader

import (
	"encoding/binary"
	"errors"
	"io"
	"time"
)

type profileHeaders struct {
	headerCount    int
	version        int
	samplingPeriod time.Duration
	is64Bit        bool
}

var errorReadingHeaders error = errors.New("Error while reading input reader")

func profileHeadersFromReader(r io.Reader) (profileHeaders, error) {
	buffer := make([]byte, 8)
	n, err := r.Read(buffer)

	// either a 32bit or 64bit reader
	var reader func(io.Reader) (uint64, error)

	if err != nil || n != 8 {
		return profileHeaders{}, errorReadingHeaders
	}

	header := profileHeaders{
		is64Bit:        true,
		samplingPeriod: time.Duration(10000) * time.Microsecond,
	}

	headerSlots := uint64(3)

	if buffer[4]|buffer[5]|buffer[6]|buffer[7] != 0 {
		header.is64Bit = false
		reader = read32BitsAsUInt64
		headerSlots = uint64(binary.LittleEndian.Uint32(buffer[4:8]))
	} else {
		reader = read64BitsAsUInt64
		header.is64Bit = true
		headerSlots, err = read64BitsAsUInt64(r)

		if err != nil {
			return profileHeaders{}, err
		}
	}

	for i := uint64(0); i < headerSlots; i++ {
		v, err := reader(r)

		if err != nil {
			return profileHeaders{}, err
		}

		switch i {
		case 0:
			header.version = int(v)
		case 1:
			header.samplingPeriod = time.Duration(v) * time.Microsecond
		}

	}

	return header, nil
}

func read32BitsAsUInt64(r io.Reader) (uint64, error) {
	var in uint32
	err := binary.Read(r, binary.LittleEndian, &in)
	return uint64(in), err
}

func read64BitsAsUInt64(r io.Reader) (uint64, error) {
	var in uint64
	err := binary.Read(r, binary.LittleEndian, &in)
	return in, err
}
