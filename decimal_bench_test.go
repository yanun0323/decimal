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

// func BenchmarkNew(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		b, _ := New(_operatorBase)
// 		_ = b
// 	}
// }

// Calculate

func BenchmarkCalculateShopSpringDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b, _ := decimal.NewFromString(_operatorBase)
		a, _ := decimal.NewFromString(_operatorAddition)
		result := b.Add(a).Sub(b)
		ss := result.String()
		_ = ss

		b.Mul(a)

	}
}

func BenchmarkCalculateDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b, _ := New(_operatorBase)
		a, _ := New(_operatorAddition)
		result := b.Add(a).Sub(a)
		_ = result
	}
}

// Shifting

func BenchmarkShiftShopSpringDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d, _ := decimal.NewFromString(_operatorBase)
		s := d.Shift(8).String()
		_ = s
		d2, _ := decimal.NewFromString(_operatorBase)
		s2 := d2.Shift(-8).String()
		_ = s2
	}
}

func BenchmarkShiftDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d, _ := New(_operatorBase)
		s := d.Shift(8)
		_ = s
		d2, _ := New(_operatorBase)
		s2 := d2.Shift(-8)
		_ = s2
	}
}

// Truncating

func BenchmarkTruncateShopSpringDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d, _ := decimal.NewFromString(_operatorBase)
		s := d.Truncate(8).String()
		_ = s
		d2, _ := decimal.NewFromString(_operatorBase)
		s2 := d2.Truncate(-8).String()
		_ = s2
	}
}

func BenchmarkTruncateDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d, _ := New(_operatorBase)
		s := d.Truncate(8)
		_ = s
		d2, _ := New(_operatorBase)
		s2 := d2.Truncate(-8)
		_ = s2
	}
}

// IsZero

func BenchmarkIsZeroShopSpringDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d1, _ := decimal.NewFromString(_operatorBase)
		_ = d1.IsZero()
	}
}

func BenchmarkIsZeroDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d1, _ := New(_operatorBase)
		_ = d1.IsZero()
	}
}

// IsPositive

func BenchmarkIsPositiveShopSpringDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d1, _ := decimal.NewFromString(_operatorBase)
		_ = d1.IsPositive()
	}
}

func BenchmarkIsPositiveDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d1, _ := New(_operatorBase)
		_ = d1.IsPositive()
	}
}

// IsNegative

func BenchmarkIsNegativeShopSpringDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d1, _ := decimal.NewFromString(_operatorBase)
		_ = d1.IsNegative()
	}
}

func BenchmarkIsNegativeDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d1, _ := New(_operatorBase)
		_ = d1.IsNegative()
	}
}

// Equal

func BenchmarkEqualShopSpringDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d1, _ := decimal.NewFromString(_operatorBase)
		d2, _ := decimal.NewFromString(_operatorAddition)
		_ = d1.Equal(d2)
	}
}

func BenchmarkEqualDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d1, _ := New(_operatorBase)
		d2, _ := New(_operatorAddition)
		_ = d1.Equal(d2)
	}
}

// Greater

func BenchmarkGreaterShopSpringDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d1, _ := decimal.NewFromString(_operatorBase)
		d2, _ := decimal.NewFromString(_operatorAddition)
		_ = d1.GreaterThan(d2)
	}
}

func BenchmarkGreaterDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d1, _ := New(_operatorBase)
		d2, _ := New(_operatorAddition)
		_ = d1.Greater(d2)
	}
}

// Less

func BenchmarkLessShopSpringDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d1, _ := decimal.NewFromString(_operatorBase)
		d2, _ := decimal.NewFromString(_operatorAddition)
		_ = d1.LessThan(d2)
	}
}

func BenchmarkLessDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d1, _ := New(_operatorBase)
		d2, _ := New(_operatorAddition)
		_ = d1.Less(d2)
	}
}

// Mul

func BenchmarkMulShopSpringDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d1, _ := decimal.NewFromString(_operatorBase)
		d2, _ := decimal.NewFromString(_operatorAddition)
		_ = d1.Mul(d2)
	}
}

func BenchmarkMulDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d1, _ := New(_operatorBase)
		d2, _ := New(_operatorAddition)
		_ = d1.Mul(d2)
	}
}
