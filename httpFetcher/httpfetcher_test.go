package httpFetcher

import (
	"bytes"
	"testing"
)

const sample_response = `num_symbols: 1
0x40103b main.LongRunningFuncion
0x4010c1 main.DoWork
0x401308 main.HttpProfileServer
0x4013e4 main.main
0x417582 runtime.main
0x419af0 runtime.goexit
0x401090 main.ShortRunningFunction
0x4010d5 main.DoWork
0x42ad82 math.Sinh
0x401036 main.LongRunningFuncion
0x401077 main.ShortRunningFunction
0x40108b main.ShortRunningFunction
0x42ae4d math.Sinh
0x42aed3 math.Sinh
0x42ad94 math.Sinh
0x42adb0 math.Sinh
0x42b4e6 math.Sqrt
0x42b4ec math.Sqrt
0x401093 main.ShortRunningFunction
0x401086 main.ShortRunningFunction
0x401081 main.ShortRunningFunction
0x42b190 math.Exp
0x42aecd math.Sinh
0x42b1a5 math.Exp
0x42b1b8 math.Exp
0x42b1ca math.Exp
0x42b1ce math.Exp
0x42b1d2 math.Exp
0x42b1e3 math.Exp
0x42b1df math.Exp
0x42b1f4 math.Exp
0x42b1fd math.Exp
0x42b20a math.Exp
0x42b217 math.Exp
0x42b213 math.Exp
0x42b220 math.Exp
0x42b224 math.Exp
0x42b22d math.Exp
0x42b231 math.Exp
0x42b23a math.Exp
0x42b23e math.Exp
0x42b247 math.Exp
0x42b24b math.Exp
0x42b254 math.Exp
0x42b258 math.Exp
0x42b265 math.Exp
0x42b276 math.Exp
0x42b283 math.Exp
0x42b287 math.Exp
0x42b294 math.Exp
0x42b298 math.Exp
0x42b2a5 math.Exp
0x42b2a9 math.Exp
0x42b2b2 math.Exp`

func TestResponseProcessing(t *testing.T) {
	functionMap, err := processSymbolResponse(bytes.NewBufferString(sample_response))

	if err != nil {
		t.Fatal(err)
	}

	expectedValue := "main.LongRunningFuncion"
	if value, present := functionMap[4198459]; !present || value != expectedValue {
		t.Errorf("Did not find expected function name for pointer, got: %v expected %v", value, expectedValue)
	}
}
