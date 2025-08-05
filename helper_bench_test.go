package decimal

import (
	"runtime"
	"testing"
)

func BenchmarkFineTuning(b *testing.B) {
	operatorBase := Require("12,345,789.00456888")
	operatorAddition := Require("789.00456888")
	beforeFunc := func() {
		operatorBase.Div(operatorAddition)
	}

	afterFunc := func() {
	}

	runtime.GC()
	runtime.GC()
	b.Run("Before", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			beforeFunc()
		}
	})

	runtime.GC()
	runtime.GC()
	b.Run("After", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			afterFunc()
		}
	})
}
