package decimal

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestDecimalFull(t *testing.T) {
	limit := decimal.RequireFromString("9999999999999999")
	start := decimal.RequireFromString("0.0000000000000000009")
	d1 := start.Copy()
	d2 := start.Copy()
	step := decimal.RequireFromString("9")
	availableGap := decimal.RequireFromString("0.0000000000000001")

	for d1.LessThanOrEqual(limit) {
		for d2.LessThanOrEqual(limit) {
			{ // add
				result := d1.Add(d2).String()
				resultDecimal := Require(d1.String()).Add(Require(d2.String())).String()
				if result != resultDecimal && decimal.RequireFromString(result).Sub(decimal.RequireFromString(resultDecimal)).Abs().GreaterThan(availableGap) {
					t.Fatalf("add mismatch, %s ~ %s, expected: %s, got: %s", d1.String(), d2.String(), result, resultDecimal)
				}
			}

			{ // add neg
				result := d1.Add(d2.Neg()).String()
				resultDecimal := Require(d1.String()).Add(Require(d2.String()).Neg()).String()
				if result != resultDecimal && decimal.RequireFromString(result).Sub(decimal.RequireFromString(resultDecimal)).Abs().GreaterThan(availableGap) {
					t.Fatalf("add neg mismatch, %s ~ %s, expected: %s, got: %s", d1.String(), d2.Neg().String(), result, resultDecimal)
				}
			}

			{ // sub
				result := d1.Sub(d2).String()
				resultDecimal := Require(d1.String()).Sub(Require(d2.String())).String()
				if result != resultDecimal && decimal.RequireFromString(result).Sub(decimal.RequireFromString(resultDecimal)).Abs().GreaterThan(availableGap) {
					t.Fatalf("sub mismatch, %s ~ %s, expected: %s, got: %s", d1.String(), d2.String(), result, resultDecimal)
				}
			}

			{ // sub neg
				result := d1.Sub(d2.Neg()).String()
				resultDecimal := Require(d1.String()).Sub(Require(d2.String()).Neg()).String()
				if result != resultDecimal && decimal.RequireFromString(result).Sub(decimal.RequireFromString(resultDecimal)).Abs().GreaterThan(availableGap) {
					t.Fatalf("sub neg mismatch, %s ~ %s, expected: %s, got: %s", d1.String(), d2.Neg().String(), result, resultDecimal)
				}
			}

			{ // mul
				result := d1.Mul(d2).String()
				resultDecimal := Require(d1.String()).Mul(Require(d2.String())).String()
				if result != resultDecimal && decimal.RequireFromString(result).Sub(decimal.RequireFromString(resultDecimal)).Abs().GreaterThan(availableGap) {
					t.Fatalf("mul mismatch, %s ~ %s, expected: %s, got: %s", d1.String(), d2.String(), result, resultDecimal)
				}
			}

			{ // mul neg
				result := d1.Mul(d2.Neg()).String()
				resultDecimal := Require(d1.String()).Mul(Require(d2.String()).Neg()).String()
				if result != resultDecimal && decimal.RequireFromString(result).Sub(decimal.RequireFromString(resultDecimal)).Abs().GreaterThan(availableGap) {
					t.Fatalf("mul neg mismatch, %s ~ %s, expected: %s, got: %s", d1.String(), d2.Neg().String(), result, resultDecimal)
				}
			}

			{ // div
				result := d1.Div(d2).String()
				resultDecimal := Require(d1.String()).Div(Require(d2.String())).String()
				if result != resultDecimal && decimal.RequireFromString(result).Sub(decimal.RequireFromString(resultDecimal)).Abs().GreaterThan(availableGap) {
					t.Fatalf("div mismatch, %s ~ %s, expected: %s, got: %s", d1.String(), d2.String(), result, resultDecimal)
				}
			}

			{ // div neg
				result := d1.Div(d2.Neg()).String()
				resultDecimal := Require(d1.String()).Div(Require(d2.String()).Neg()).String()
				if result != resultDecimal && decimal.RequireFromString(result).Sub(decimal.RequireFromString(resultDecimal)).Abs().GreaterThan(availableGap) {
					t.Fatalf("div neg mismatch, %s ~ %s, expected: %s, got: %s", d1.String(), d2.Neg().String(), result, resultDecimal)
				}
			}
			d2 = d2.Mul(step).Add(d2)
		}

		d1 = d1.Mul(step).Add(d1)
	}
}

func TestDecimalPow(t *testing.T) {
	d1 := decimal.RequireFromString("123456")
	d2 := decimal.RequireFromString("56789")
	result := d1.Pow(d2).String()
	resultDecimal := decimal.RequireFromString(d1.String()).Pow(decimal.RequireFromString(d2.String())).String()
	if result != resultDecimal {
		t.Fatalf("pow mismatch, %s ~ %s, expected: %s, got: %s", d1.String(), d2.String(), result, resultDecimal)
	}
}
