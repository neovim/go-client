package plugin

import "testing"

type env struct {
	GOROOT string `eval:"$GOROOT"`
	GOPATH string `eval:"$GOPATH"`
}

var evalTests = []struct {
	fn   interface{}
	eval string
}{
	{func(x *struct {
		X  int    `eval:"1"`
		YY string `eval:"'hello'" msgpack:"Y"`
		Z  int
	}) {
	}, `{'X': 1, 'Y': 'hello'}`},

	// Nested struct
	{func(x *struct {
		Env env
	}) {
	}, `{'Env': {'GOROOT': $GOROOT, 'GOPATH': $GOPATH}}`},
}

func TestEval(t *testing.T) {
	for _, tt := range evalTests {
		eval := eval("*", tt.fn)
		if eval != tt.eval {
			t.Errorf("eval(%T) returned %q, want %q", tt.fn, eval, tt.eval)
		}
	}
}
