package decimal

import "testing"

var (
	benchD1  = mustDecimalBench("123456789.987654321")
	benchD2  = mustDecimalBench("-987654321.123456789")
	benchD3  = NewDecimal256FromInt(2)
	benchBin = func() []byte {
		b, _ := benchD1.EncodeBinary()
		return b
	}()
	benchJSON = func() []byte {
		b, _ := benchD1.EncodeJSON()
		return b
	}()
	benchBuf = make([]byte, 0, 256)
)

var (
	benchSinkDecimal Decimal256
	benchSinkBool    bool
	benchSinkInt     int
	benchSinkInt64   int64
	benchSinkFloat   float64
	benchSinkString  string
	benchSinkBytes   []byte
	benchSinkErr     error
)

func mustDecimalBench(s string) Decimal256 {
	d, err := NewDecimal256FromString(s)
	if err != nil {
		panic(err)
	}
	return d
}

func BenchmarkDecimal256_Constructors(b *testing.B) {
	b.Run("NewDecimal256", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = NewDecimal256(123456789, 987654321)
		}
	})
	b.Run("NewDecimal256FromString", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal, benchSinkErr = NewDecimal256FromString("123456789.987654321")
		}
	})
	b.Run("NewDecimal256FromInt", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = NewDecimal256FromInt(123456789)
		}
	})
	b.Run("NewDecimal256FromFloat", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal, benchSinkErr = NewDecimal256FromFloat(12345.5)
		}
	})
}

func BenchmarkDecimal256_Conversions(b *testing.B) {
	b.Run("Int64", func(b *testing.B) {
		for b.Loop() {
			ip, fp := benchD1.Int64()
			benchSinkInt64 = ip + fp
		}
	})
	b.Run("Float64", func(b *testing.B) {
		for b.Loop() {
			benchSinkFloat = benchD1.Float64()
		}
	})
	b.Run("String", func(b *testing.B) {
		for b.Loop() {
			benchSinkString = benchD1.String()
		}
	})
	b.Run("StringFixed", func(b *testing.B) {
		for b.Loop() {
			benchSinkString = benchD1.StringFixed(16)
		}
	})
}

func BenchmarkDecimal256_Checking(b *testing.B) {
	b.Run("IsZero", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD1.IsZero()
		}
	})
	b.Run("IsPositive", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD1.IsPositive()
		}
	})
	b.Run("IsNegative", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD2.IsNegative()
		}
	})
	b.Run("Sign", func(b *testing.B) {
		for b.Loop() {
			benchSinkInt = benchD2.Sign()
		}
	})
}

func BenchmarkDecimal256_Modification(b *testing.B) {
	b.Run("Neg", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.Neg()
		}
	})
	b.Run("Inv", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.Inv()
		}
	})
	b.Run("Abs", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD2.Abs()
		}
	})
	b.Run("Truncate", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.Truncate(8)
		}
	})
	b.Run("Shift", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.Shift(3)
		}
	})
	b.Run("Round", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.Round(8)
		}
	})
	b.Run("RoundAwayFromZero", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.RoundAwayFromZero(8)
		}
	})
	b.Run("RoundTowardToZero", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.RoundTowardToZero(8)
		}
	})
	b.Run("Ceil", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.Ceil(2)
		}
	})
	b.Run("Floor", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.Floor(2)
		}
	})
}

func BenchmarkDecimal256_Comparison(b *testing.B) {
	b.Run("Equal", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD1.Equal(benchD1)
		}
	})
	b.Run("GreaterThan", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD1.GreaterThan(benchD2)
		}
	})
	b.Run("LessThan", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD2.LessThan(benchD1)
		}
	})
	b.Run("GreaterOrEqual", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD1.GreaterOrEqual(benchD1)
		}
	})
	b.Run("LessOrEqual", func(b *testing.B) {
		for b.Loop() {
			benchSinkBool = benchD2.LessOrEqual(benchD2)
		}
	})
}

func BenchmarkDecimal256_Arithmetic(b *testing.B) {
	b.Run("Add", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.Add(benchD2)
		}
	})
	b.Run("Sub", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.Sub(benchD2)
		}
	})
	b.Run("Mul", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.Mul(benchD3)
		}
	})
	b.Run("Div", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.Div(benchD3)
		}
	})
	b.Run("Mod", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.Mod(benchD3)
		}
	})
}

func BenchmarkDecimal256_Transcendental(b *testing.B) {
	b.Run("Pow", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD3.Pow(NewDecimal256FromInt(7))
		}
	})
	b.Run("Sqrt", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.Abs().Sqrt()
		}
	})
	b.Run("Exp", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = NewDecimal256FromInt(1).Exp()
		}
	})
	b.Run("Log", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.Abs().Log()
		}
	})
	b.Run("Log2", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.Abs().Log2()
		}
	})
	b.Run("Log10", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal = benchD1.Abs().Log10()
		}
	})
}

func BenchmarkDecimal256_EncodeDecode(b *testing.B) {
	b.Run("EncodeBinary", func(b *testing.B) {
		for b.Loop() {
			benchSinkBytes, benchSinkErr = benchD1.EncodeBinary()
		}
	})
	b.Run("AppendBinary", func(b *testing.B) {
		for b.Loop() {
			dst := benchD1.AppendBinary(benchBuf[:0])
			benchSinkInt = len(dst)
		}
	})
	b.Run("NewFromBinary", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal, benchSinkErr = NewDecimal256FromBinary(benchBin)
		}
	})
	b.Run("EncodeJSON", func(b *testing.B) {
		for b.Loop() {
			benchSinkBytes, benchSinkErr = benchD1.EncodeJSON()
		}
	})
	b.Run("AppendJSON", func(b *testing.B) {
		for b.Loop() {
			dst := benchD1.AppendJSON(benchBuf[:0])
			benchSinkInt = len(dst)
		}
	})
	b.Run("NewFromJSON", func(b *testing.B) {
		for b.Loop() {
			benchSinkDecimal, benchSinkErr = NewDecimal256FromJSON(benchJSON)
		}
	})
	b.Run("AppendString", func(b *testing.B) {
		for b.Loop() {
			dst := benchD1.AppendString(benchBuf[:0])
			benchSinkInt = len(dst)
		}
	})
	b.Run("AppendStringFixed", func(b *testing.B) {
		for b.Loop() {
			dst := benchD1.AppendStringFixed(benchBuf[:0], 16)
			benchSinkInt = len(dst)
		}
	})
}
