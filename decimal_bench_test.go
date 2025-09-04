package decimal

import (
	"runtime"
	"testing"

	"github.com/shopspring/decimal"
)

const (
	_operatorBase     = "12,345,789.00456888"
	_operatorAddition = "789.00456888"
)

func Run(b *testing.B, shop, dec func(b *testing.B)) {

	runtime.GC()
	runtime.GC()

	b.Run("ShopSpringDecimal", func(b *testing.B) {
		shop(b)
	})

	runtime.GC()
	runtime.GC()

	b.Run("Decimal", func(b *testing.B) {
		dec(b)
	})
}

func BenchmarkNew(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for i := 0; i < 1; i++ {
					dd, _ := decimal.NewFromString(_operatorBase)
					_ = dd
				}
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for i := 0; i < 1; i++ {
					dd, _ := New(_operatorBase)
					_ = dd
				}
			}
		},
	)
}

func BenchmarkNewFromFloat(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				dd := decimal.NewFromFloat(123456789.00456888)
				_ = dd
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				dd := NewFromFloat(123456789.00456888)
				_ = dd
			}
		},
	)
}

func BenchmarkNewFromInt(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				dd := decimal.NewFromInt(123456789)
				_ = dd
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				dd := NewFromInt(123456789)
				_ = dd
			}
		},
	)
}

func BenchmarkStringFixed(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := decimal.NewFromString(_operatorBase)
				_ = d1.StringFixed(2)
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := New(_operatorBase)
				_ = d1.StringFixed(2)
			}
		},
	)
}

func BenchmarkAbs(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := decimal.NewFromString(_operatorBase)
				_ = d1.Abs()
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := New(_operatorBase)
				_ = d1.Abs()
			}
		},
	)
}

func BenchmarkNeg(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := decimal.NewFromString(_operatorBase)
				_ = d1.Neg()
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := New(_operatorBase)
				_ = d1.Neg()
			}
		},
	)
}

func BenchmarkTruncate(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d, _ := decimal.NewFromString(_operatorBase)
				s := d.Truncate(8).String()
				_ = s
				d2, _ := decimal.NewFromString(_operatorBase)
				s2 := d2.Truncate(-8).String()
				_ = s2
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d, _ := New(_operatorBase)
				s := d.Truncate(8)
				_ = s
				d2, _ := New(_operatorBase)
				s2 := d2.Truncate(-8)
				_ = s2
			}
		},
	)
}

func BenchmarkCeil(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d, _ := decimal.NewFromString(_operatorBase)
				s := d.RoundCeil(8).String()
				_ = s
				d2, _ := decimal.NewFromString(_operatorBase)
				s2 := d2.RoundCeil(-8).String()
				_ = s2
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d, _ := New(_operatorBase)
				s := d.Ceil(8)
				_ = s
				d2, _ := New(_operatorBase)
				s2 := d2.Ceil(-8)
				_ = s2
			}
		},
	)
}

func BenchmarkRound(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d, _ := decimal.NewFromString(_operatorBase)
				s := d.Round(8).String()
				_ = s
				d2, _ := decimal.NewFromString(_operatorBase)
				s2 := d2.Round(-8).String()
				_ = s2
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d, _ := New(_operatorBase)
				s := d.Round(8)
				_ = s
				d2, _ := New(_operatorBase)
				s2 := d2.Round(-8)
				_ = s2
			}
		},
	)
}

func BenchmarkRoundAwayFromZero(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d, _ := decimal.NewFromString(_operatorBase)
				s := d.RoundUp(8).String()
				_ = s
				d2, _ := decimal.NewFromString(_operatorBase)
				s2 := d2.RoundUp(-8).String()
				_ = s2
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d, _ := New(_operatorBase)
				s := d.RoundAwayFromZero(8)
				_ = s
				d2, _ := New(_operatorBase)
				s2 := d2.RoundAwayFromZero(-8)
				_ = s2
			}
		},
	)
}

func BenchmarkRoundTowardToZero(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d, _ := decimal.NewFromString(_operatorBase)
				s := d.RoundDown(8).String()
				_ = s
				d2, _ := decimal.NewFromString(_operatorBase)
				s2 := d2.RoundDown(-8).String()
				_ = s2
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d, _ := New(_operatorBase)
				s := d.RoundTowardToZero(8)
				_ = s
				d2, _ := New(_operatorBase)
				s2 := d2.RoundTowardToZero(-8)
				_ = s2
			}
		},
	)
}

func BenchmarkFloor(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d, _ := decimal.NewFromString(_operatorBase)
				s := d.RoundFloor(8).String()
				_ = s
				d2, _ := decimal.NewFromString(_operatorBase)
				s2 := d2.RoundFloor(-8).String()
				_ = s2
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d, _ := New(_operatorBase)
				s := d.Floor(8)
				_ = s
				d2, _ := New(_operatorBase)
				s2 := d2.Floor(-8)
				_ = s2
			}
		},
	)
}

func BenchmarkShift(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d, _ := decimal.NewFromString(_operatorBase)
				s := d.Shift(8).String()
				_ = s
				d2, _ := decimal.NewFromString(_operatorBase)
				s2 := d2.Shift(-8).String()
				_ = s2
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d, _ := New(_operatorBase)
				s := d.Shift(8)
				_ = s
				d2, _ := New(_operatorBase)
				s2 := d2.Shift(-8)
				_ = s2
			}
		},
	)
}

func BenchmarkAdd(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b, _ := decimal.NewFromString(_operatorBase)
				a, _ := decimal.NewFromString(_operatorAddition)
				result := b.Add(a).Add(b)
				ss := result.String()
				_ = ss
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b, _ := New(_operatorBase)
				a, _ := New(_operatorAddition)
				result := b.Add(a).Add(a)
				_ = result
			}
		},
	)
}

func BenchmarkSub(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b, _ := decimal.NewFromString(_operatorBase)
				a, _ := decimal.NewFromString(_operatorAddition)
				result := b.Sub(a).Sub(b)
				ss := result.String()
				_ = ss
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b, _ := New(_operatorBase)
				a, _ := New(_operatorAddition)
				result := b.Sub(a).Sub(a)
				_ = result
			}
		},
	)
}

func BenchmarkIsZero(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := decimal.NewFromString(_operatorBase)
				_ = d1.IsZero()
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := New(_operatorBase)
				_ = d1.IsZero()
			}
		},
	)
}

func BenchmarkIsPositive(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := decimal.NewFromString(_operatorBase)
				_ = d1.IsPositive()
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := New(_operatorBase)
				_ = d1.IsPositive()
			}
		},
	)
}

func BenchmarkIsNegative(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := decimal.NewFromString(_operatorBase)
				_ = d1.IsNegative()
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := New(_operatorBase)
				_ = d1.IsNegative()
			}
		},
	)
}

func BenchmarkCmp(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := decimal.NewFromString(_operatorBase)
				d2, _ := decimal.NewFromString(_operatorAddition)
				_ = d1.Cmp(d2)
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := New(_operatorBase)
				d2, _ := New(_operatorAddition)
				_ = d1.Cmp(d2)
			}
		},
	)
}

func BenchmarkEqual(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := decimal.NewFromString(_operatorBase)
				d2, _ := decimal.NewFromString(_operatorAddition)
				_ = d1.Equal(d2)
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := New(_operatorBase)
				d2, _ := New(_operatorAddition)
				_ = d1.Equal(d2)
			}
		},
	)
}

func BenchmarkGreater(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := decimal.NewFromString(_operatorBase)
				d2, _ := decimal.NewFromString(_operatorAddition)
				_ = d1.GreaterThan(d2)
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := New(_operatorBase)
				d2, _ := New(_operatorAddition)
				_ = d1.Greater(d2)
			}
		},
	)
}

func BenchmarkLess(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := decimal.NewFromString(_operatorBase)
				d2, _ := decimal.NewFromString(_operatorAddition)
				_ = d1.LessThan(d2)
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := New(_operatorBase)
				d2, _ := New(_operatorAddition)
				_ = d1.Less(d2)
			}
		},
	)
}

func BenchmarkGreaterOrEqual(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := decimal.NewFromString(_operatorBase)
				d2, _ := decimal.NewFromString(_operatorAddition)
				_ = d1.GreaterThanOrEqual(d2)
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := New(_operatorBase)
				d2, _ := New(_operatorAddition)
				_ = d1.GreaterOrEqual(d2)
			}
		},
	)
}

func BenchmarkLessOrEqual(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := decimal.NewFromString(_operatorBase)
				d2, _ := decimal.NewFromString(_operatorAddition)
				_ = d1.LessThanOrEqual(d2)
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := New(_operatorBase)
				d2, _ := New(_operatorAddition)
				_ = d1.LessOrEqual(d2)
			}
		},
	)
}

func BenchmarkSign(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := decimal.NewFromString(_operatorBase)
				s := d1.Sign()
				_ = s + 2
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := New(_operatorBase)
				s := d1.Sign()
				_ = s + 2
			}
		},
	)
}

func BenchmarkMul(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := decimal.NewFromString(_operatorBase)
				d2, _ := decimal.NewFromString(_operatorAddition)
				_ = d1.Mul(d2)
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := New(_operatorBase)
				d2, _ := New(_operatorAddition)
				_ = d1.Mul(d2)
			}
		},
	)
}

func BenchmarkDiv(b *testing.B) {
	Run(b,
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := decimal.NewFromString(_operatorBase)
				d2, _ := decimal.NewFromString(_operatorAddition)
				_ = d1.Div(d2)
			}
		},
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d1, _ := New(_operatorBase)
				d2, _ := New(_operatorAddition)
				_ = d1.Div(d2)
			}
		},
	)
}
