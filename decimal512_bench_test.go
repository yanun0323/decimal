package decimal

import "testing"

var (
	benchD512   = mustDecimal512Bench("123456789.987654321")
	benchD512b  = mustDecimal512Bench("-987654321.123456789")
	benchD512c  = NewDecimal512FromInt(2)
	benchBin512 = func() []byte {
		b, _ := benchD512.EncodeBinary()
		return b
	}()
	benchJSON512 = func() []byte {
		b, _ := benchD512.EncodeJSON()
		return b
	}()
	benchBuf512 = make([]byte, 0, 256)
)

var (
	benchSinkDecimal512 Decimal512
)

func mustDecimal512Bench(s string) Decimal512 {
	d, err := NewDecimal512FromString(s)
	if err != nil {
		panic(err)
	}
	return d
}

func BenchmarkDecimal512_Constructors(b *testing.B) {
	b.Run("NewDecimal512", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = NewDecimal512(123456789, 987654321)
		}
	})
	b.Run("NewDecimal512FromString", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512, benchSinkErr = NewDecimal512FromString("123456789.987654321")
		}
	})
	b.Run("NewDecimal512FromInt", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = NewDecimal512FromInt(123456789)
		}
	})
	b.Run("NewDecimal512FromFloat", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512, benchSinkErr = NewDecimal512FromFloat(12345.5)
		}
	})
}

func BenchmarkDecimal512_Conversions(b *testing.B) {
	b.Run("Int64", func(b *testing.B) {
		for b.Loop() {
			ip, fp := benchD512.Int64()
			benchSinkInt64 = ip + fp
		}
	})
	b.Run("Float64", func(b *testing.B) {
		for b.Loop() {
			benchSinkFloat = benchD512.Float64()
		}
	})
	b.Run("String", func(b *testing.B) {
		for b.Loop() {
			benchSinkString = benchD512.String()
		}
	})
	b.Run("StringFixed", func(b *testing.B) {
		for b.Loop() {
			benchSinkString = benchD512.StringFixed(16)
		}
	})
}

func BenchmarkDecimal512_Checking(b *testing.B) {
	b.Run("IsZero", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD512.IsZero()
		}
	})
	b.Run("IsPositive", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD512.IsPositive()
		}
	})
	b.Run("IsNegative", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD512b.IsNegative()
		}
	})
	b.Run("Sign", func(b *testing.B) {
		for b.Loop() {
			benchSinkInt = benchD512b.Sign()
		}
	})
}

func BenchmarkDecimal512_Modification(b *testing.B) {
	b.Run("Neg", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Neg()
		}
	})
	b.Run("Inv", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Inv()
		}
	})
	b.Run("Abs", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512b.Abs()
		}
	})
	b.Run("Truncate", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Truncate(8)
		}
	})
	b.Run("Shift", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Shift(3)
		}
	})
	b.Run("Round", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Round(8)
		}
	})
	b.Run("RoundAwayFromZero", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.RoundAwayFromZero(8)
		}
	})
	b.Run("RoundTowardToZero", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.RoundTowardToZero(8)
		}
	})
	b.Run("Ceil", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Ceil(2)
		}
	})
	b.Run("Floor", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Floor(2)
		}
	})
}

func BenchmarkDecimal512_Comparison(b *testing.B) {
	b.Run("Equal", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD512.Equal(benchD512)
		}
	})
	b.Run("GreaterThan", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD512.GreaterThan(benchD512b)
		}
	})
	b.Run("LessThan", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD512b.LessThan(benchD512)
		}
	})
	b.Run("GreaterOrEqual", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD512.GreaterOrEqual(benchD512)
		}
	})
	b.Run("LessOrEqual", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD512b.LessOrEqual(benchD512b)
		}
	})
}

func BenchmarkDecimal512_Arithmetic(b *testing.B) {
	b.Run("Add", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Add(benchD512b)
		}
	})
	b.Run("Sub", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Sub(benchD512b)
		}
	})
	b.Run("Mul", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Mul(benchD512b)
		}
	})
	b.Run("Div", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Div(benchD512c)
		}
	})
	b.Run("Mod", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Mod(benchD512c)
		}
	})
}

func BenchmarkDecimal512_Transcendental(b *testing.B) {
	b.Run("Pow", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Pow(benchD512c)
		}
	})
	b.Run("Sqrt", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Sqrt()
		}
	})
	b.Run("Exp", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Exp()
		}
	})
	b.Run("Log", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Log()
		}
	})
	b.Run("Log2", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Log2()
		}
	})
	b.Run("Log10", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512 = benchD512.Log10()
		}
	})
}

func BenchmarkDecimal512_EncodeDecode(b *testing.B) {
	b.Run("EncodeBinary", func(b *testing.B) {
		for b.Loop() {
			benchSinkBytes, benchSinkErr = benchD512.EncodeBinary()
		}
	})
	b.Run("AppendBinary", func(b *testing.B) {
		for b.Loop() {
			benchBuf512 = benchD512.AppendBinary(benchBuf512[:0])
			benchSinkInt = len(benchBuf512)
		}
	})
	b.Run("NewFromBinary", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512, benchSinkErr = NewDecimal512FromBinary(benchBin512)
		}
	})
	b.Run("EncodeJSON", func(b *testing.B) {
		for b.Loop() {
			benchSinkBytes, benchSinkErr = benchD512.EncodeJSON()
		}
	})
	b.Run("AppendJSON", func(b *testing.B) {
		for b.Loop() {
			benchBuf512 = benchD512.AppendJSON(benchBuf512[:0])
			benchSinkInt = len(benchBuf512)
		}
	})
	b.Run("NewFromJSON", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal512, benchSinkErr = NewDecimal512FromJSON(benchJSON512)
		}
	})
	b.Run("AppendString", func(b *testing.B) {
		for b.Loop() {
			benchBuf512 = benchD512.AppendString(benchBuf512[:0])
			benchSinkInt = len(benchBuf512)
		}
	})
	b.Run("AppendStringFixed", func(b *testing.B) {
		for b.Loop() {
			benchBuf512 = benchD512.AppendStringFixed(benchBuf512[:0], 16)
			benchSinkInt = len(benchBuf512)
		}
	})
}
