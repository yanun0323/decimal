package decimal

import (
	"testing"

	"github.com/shopspring/decimal"
)

// command:
// go test -bench=. -run=none -benchmem .
const (
	_operatorBase     = "12,345,789.00456888"
	_operatorAddition = "789.00456888"
)

func BenchmarkNew(b *testing.B) {
	b.SkipNow()
	for i := 0; i < b.N; i++ {
		b, _ := New(_operatorBase)
		_ = b
	}
}

func Run(b *testing.B, shop, dec func(b *testing.B)) {
	b.Run("ShopSpringDecimal", func(b *testing.B) {
		shop(b)
	})

	b.Run("Decimal", func(b *testing.B) {
		dec(b)
	})
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
