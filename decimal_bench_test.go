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
