package plugin

import (
	"testing"
)

var evalString string

func BenchmarkEval(b *testing.B) {
	fn := func(x *struct {
		X  int    `eval:"1"`
		YY string `eval:"'hello'" msgpack:"Y"`
		Z  int
	}) {
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		evalString = eval("*", fn)
	}
}
