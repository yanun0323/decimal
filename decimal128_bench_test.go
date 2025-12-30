package decimal

import "testing"

var (
	benchD128   = mustDecimal128Bench("123456.987654321")
	benchD128b  = mustDecimal128Bench("-987654.123456789")
	benchD128c  = New128FromInt(2)
	benchBin128 = func() []byte {
		b, _ := benchD128.EncodeBinary()
		return b
	}()
	benchJSON128 = func() []byte {
		b, _ := benchD128.EncodeJSON()
		return b
	}()
	benchBuf128 = make([]byte, 0, 256)
)

var (
	benchSinkDecimal128 Decimal128
)

func mustDecimal128Bench(s string) Decimal128 {
	d, err := New128FromString(s)
	if err != nil {
		panic(err)
	}
	return d
}

func BenchmarkDecimal128_Constructors(b *testing.B) {
	b.Run("New128", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = New128(123456, 987654321)
		}
	})
	b.Run("New128FromString", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128, benchSinkErr = New128FromString("123456.987654321")
		}
	})
	b.Run("New128FromInt", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = New128FromInt(123456)
		}
	})
	b.Run("New128FromFloat", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128, benchSinkErr = New128FromFloat(12345.5)
		}
	})
}

func BenchmarkDecimal128_Conversions(b *testing.B) {
	b.Run("Int64", func(b *testing.B) {
		for b.Loop() {
			ip, fp := benchD128.Int64()
			benchSinkInt64 = ip + fp
		}
	})
	b.Run("Float64", func(b *testing.B) {
		for b.Loop() {
			benchSinkFloat = benchD128.Float64()
		}
	})
	b.Run("String", func(b *testing.B) {
		for b.Loop() {
			benchSinkString = benchD128.String()
		}
	})
	b.Run("StringFixed", func(b *testing.B) {
		for b.Loop() {
			benchSinkString = benchD128.StringFixed(16)
		}
	})
}

func BenchmarkDecimal128_Checking(b *testing.B) {
	b.Run("IsZero", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD128.IsZero()
		}
	})
	b.Run("IsPositive", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD128.IsPositive()
		}
	})
	b.Run("IsNegative", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD128b.IsNegative()
		}
	})
	b.Run("Sign", func(b *testing.B) {
		for b.Loop() {
			benchSinkInt = benchD128b.Sign()
		}
	})
}

func BenchmarkDecimal128_Modification(b *testing.B) {
	b.Run("Neg", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Neg()
		}
	})
	b.Run("Inv", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Inv()
		}
	})
	b.Run("Abs", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128b.Abs()
		}
	})
	b.Run("Truncate", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Truncate(8)
		}
	})
	b.Run("Shift", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Shift(3)
		}
	})
	b.Run("Round", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Round(8)
		}
	})
	b.Run("RoundAwayFromZero", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.RoundAwayFromZero(8)
		}
	})
	b.Run("RoundTowardToZero", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.RoundTowardToZero(8)
		}
	})
	b.Run("Ceil", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Ceil(2)
		}
	})
	b.Run("Floor", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Floor(2)
		}
	})
}

func BenchmarkDecimal128_Comparison(b *testing.B) {
	b.Run("Equal", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD128.Equal(benchD128)
		}
	})
	b.Run("GreaterThan", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD128.GreaterThan(benchD128b)
		}
	})
	b.Run("LessThan", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD128b.LessThan(benchD128)
		}
	})
	b.Run("GreaterOrEqual", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD128.GreaterOrEqual(benchD128)
		}
	})
	b.Run("LessOrEqual", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD128b.LessOrEqual(benchD128b)
		}
	})
}

func BenchmarkDecimal128_Arithmetic(b *testing.B) {
	b.Run("Add", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Add(benchD128b)
		}
	})
	b.Run("Sub", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Sub(benchD128b)
		}
	})
	b.Run("Mul", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Mul(benchD128b)
		}
	})
	b.Run("Div", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Div(benchD128c)
		}
	})
	b.Run("Mod", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Mod(benchD128c)
		}
	})
}

func BenchmarkDecimal128_Transcendental(b *testing.B) {
	b.Run("Pow", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Pow(benchD128c)
		}
	})
	b.Run("Sqrt", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Sqrt()
		}
	})
	b.Run("Exp", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Exp()
		}
	})
	b.Run("Log", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Log()
		}
	})
	b.Run("Log2", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Log2()
		}
	})
	b.Run("Log10", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128 = benchD128.Log10()
		}
	})
}

func BenchmarkDecimal128_EncodeDecode(b *testing.B) {
	b.Run("EncodeBinary", func(b *testing.B) {
		for b.Loop() {
			benchSinkBytes, benchSinkErr = benchD128.EncodeBinary()
		}
	})
	b.Run("AppendBinary", func(b *testing.B) {
		for b.Loop() {
			benchBuf128 = benchD128.AppendBinary(benchBuf128[:0])
			benchSinkInt = len(benchBuf128)
		}
	})
	b.Run("NewFromBinary", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128, benchSinkErr = New128FromBinary(benchBin128)
		}
	})
	b.Run("EncodeJSON", func(b *testing.B) {
		for b.Loop() {
			benchSinkBytes, benchSinkErr = benchD128.EncodeJSON()
		}
	})
	b.Run("AppendJSON", func(b *testing.B) {
		for b.Loop() {
			benchBuf128 = benchD128.AppendJSON(benchBuf128[:0])
			benchSinkInt = len(benchBuf128)
		}
	})
	b.Run("NewFromJSON", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal128, benchSinkErr = New128FromJSON(benchJSON128)
		}
	})
	b.Run("AppendString", func(b *testing.B) {
		for b.Loop() {
			benchBuf128 = benchD128.AppendString(benchBuf128[:0])
			benchSinkInt = len(benchBuf128)
		}
	})
	b.Run("AppendStringFixed", func(b *testing.B) {
		for b.Loop() {
			benchBuf128 = benchD128.AppendStringFixed(benchBuf128[:0], 16)
			benchSinkInt = len(benchBuf128)
		}
	})
}
