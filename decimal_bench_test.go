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

func Benchmark_NewDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b, _ := NewDecimal(_operatorBase)
		_ = b
	}
}

// Calculate

func Benchmark_Calculate_Decimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b, _ := decimal.NewFromString(_operatorBase)
		a, _ := decimal.NewFromString(_operatorAddition)
		result := b.Add(a).Sub(b)
		ss := result.String()
		_ = ss

	}
}

func Benchmark_Calculate_StringDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b, _ := NewDecimal(_operatorBase)
		a, _ := NewDecimal(_operatorAddition)
		result := b.Add(a).Sub(a)
		_ = result
	}
}

// Shifting

// TODO: Shift Decimal

func Benchmark_Shift_StringDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d, _ := NewDecimal(_operatorBase)
		s := d.Shift(8)
		_ = s
		d2, _ := NewDecimal(_operatorBase)
		s2 := d2.Shift(-8)
		_ = s2
	}
}
