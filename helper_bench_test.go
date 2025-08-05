package decimal

import (
	"runtime"
	"testing"
)

func BenchmarkFineTuning(b *testing.B) {
	beforeFunc := func() {
		slice := make([]byte, 1000)
		slice2 := make([]byte, 1001)[:1000]
		insert(slice, 100, '0')
		insert(slice2, 100, '0')
	}

	afterFunc := func() {
		slice := make([]byte, 1000)
		slice2 := make([]byte, 1001)[:1000]
		insert(slice, 100, '0')
		insert(slice2, 100, '0')
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
