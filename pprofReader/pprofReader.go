/*
This package is designed to read the pprof format as described:
http://google-perftools.googlecode.com/svn/trunk/doc/cpuprofile-fileformat.html


*/
package pprofReader

import (
	"io"
)

type Profile struct {
	headers profileHeaders
	Samples ProfileSamples
	// TODO: need mapped objects? probably not for Go
}

// Read a profile from disk or from a http response body
func ReadProfile(r io.Reader) (p *Profile, err error) {
	p = &Profile{}
	p.headers, err = profileHeadersFromReader(r)

	if err != nil {
		return
	}

	var reader func(io.Reader) (Sample, error)

	if p.headers.is64Bit {
		reader = read64BitSample
	} else {
		reader = read32BitSample
	}

	for {
		var s Sample
		s, err = reader(r)
		if err != nil {
			return nil, err
		}
		if !s.IsEnd() {
			p.Samples = append(p.Samples, s)
		} else {
			break
		}
	}

	return p, err
}
