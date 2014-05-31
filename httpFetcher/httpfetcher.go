/*
This package implements an interface to interact with the http server behind net/http/pprof
We don't need to worry too much about the implementation on the backend, just the interface to it

TODO worry about timeouts on long profiles
TODO allow the http fetcher to run profiles for different lengths of time
*/
package httpFetcher

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Fetcher interface {
	// obtain a readable profile
	GetProfile() (io.Reader, error)

	// obtain information about a given function pointer
	GetFunctionNames(fnPtrs []uint64) (map[uint64]string, error)
}

// Fetches Pprof results from the standard golang pprof server
type HttpFetcher struct {
	// default: localhost
	Server string
	// default: 6060
	Port uint
	// default: 30s
	ProfileDuration time.Duration
}

func (f *HttpFetcher) getVars() (server string, port uint, duration time.Duration) {
	server = "localhost"
	if f.Server != "" {
		server = f.Server
	}

	port = 6060
	if f.Port != 0 {
		port = f.Port
	}

	duration = time.Second * time.Duration(30)
	if f.ProfileDuration != time.Duration(0) {
		duration = f.ProfileDuration
	}

	return
}

// Fetch a new profile from the http server
func (f *HttpFetcher) GetProfile() (io.Reader, error) {
	server, port, _ := f.getVars()

	urlStr := fmt.Sprintf("http://%v:%v/debug/pprof/profile", server, port)
	s := uint(f.ProfileDuration.Seconds())
	response, err := http.PostForm(urlStr, url.Values{"seconds": []string{fmt.Sprint(s)}})

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	buffer := bytes.NewBuffer(nil)
	io.Copy(buffer, response.Body)

	return buffer, nil
}

// Fetch symbols from /debug/pprof/symbol
// Use a POST request, with function pointers separated by + symbols
// the response uses the following formatting - "%#x %s\n" for each function pointer and name
// (see: Symbol in pprof.go)
func (f *HttpFetcher) GetFunctionNames(fnPtrs []uint64) (map[uint64]string, error) {
	server, port, _ := f.getVars()

	ptrsAsStrings := make([]string, len(fnPtrs))
	for i, v := range fnPtrs {
		ptrsAsStrings[i] = strconv.FormatUint(v, 10)
	}

	bodyString := bytes.NewBufferString(strings.Join(ptrsAsStrings, "+"))

	response, err := http.Post(
		fmt.Sprintf("http://%v:%v/debug/pprof/symbol", server, port),
		"text/plain; charset=utf-8",
		bodyString,
	)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return processSymbolResponse(response.Body)
}

func processSymbolResponse(r io.Reader) (result map[uint64]string, err error) {
	b := bufio.NewReader(r)
	var line string

	line, err = b.ReadString('\n')
	if err != nil {
		return
	}

	if line != "num_symbols: 1\n" {
		err = fmt.Errorf("Expected \"num_symbols: 1\", got: \"%v\"\n", line)
		return
	}

	result = make(map[uint64]string)

	for {
		line, err = b.ReadString(' ')
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}

		functionPtrString := line[:len(line)-1]
		var i uint64
		i, err = strconv.ParseUint(functionPtrString, 0, 64)

		if err != nil {
			return
		}

		line, err = b.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}

		functionName := line[:len(line)-1]

		result[i] = functionName

	}

	return
}
