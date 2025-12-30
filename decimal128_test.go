package decimal

import (
	"bytes"
	"math"
	"testing"
)

func mustDecimal128(t *testing.T, s string) Decimal128 {
	t.Helper()
	d, err := NewDecimal128FromString(s)
	if err != nil {
		t.Fatalf("NewDecimal128FromString(%q) error: %v", s, err)
	}
	return d
}

func TestDecimal128Constructors(t *testing.T) {
	d1 := NewDecimal128(123456, 987654321)
	if got := d1.String(); got != "123456.987654321" {
		t.Fatalf("NewDecimal128 string mismatch: %s", got)
	}

	d1b := NewDecimal128(1, 1234567890123456789)
	if got := d1b.String(); got != "1.1234567890123456" {
		t.Fatalf("NewDecimal128 truncation mismatch: %s", got)
	}

	d2 := NewDecimal128FromInt(-123456)
	if got := d2.String(); got != "-123456" {
		t.Fatalf("NewDecimal128FromInt string mismatch: %s", got)
	}

	d2b := NewDecimal128FromInt(10000000000000000)
	if got := d2b.String(); got != "0" {
		t.Fatalf("NewDecimal128FromInt truncation mismatch: %s", got)
	}

	d3, err := NewDecimal128FromFloat(1.25)
	if err != nil {
		t.Fatalf("NewDecimal128FromFloat error: %v", err)
	}
	if diff := math.Abs(d3.Float64() - 1.25); diff > 1e-9 {
		t.Fatalf("NewDecimal128FromFloat value mismatch: diff=%g", diff)
	}

	d4, err := NewDecimal128FromString("  +1_234.4500e0 ")
	if err != nil {
		t.Fatalf("NewDecimal128FromString error: %v", err)
	}
	if got := d4.String(); got != "1234.45" {
		t.Fatalf("NewDecimal128FromString string mismatch: %s", got)
	}

	d5, err := NewDecimal128FromString("-.5")
	if err != nil {
		t.Fatalf("NewDecimal128FromString(-.5) error: %v", err)
	}
	if got := d5.String(); got != "-0.5" {
		t.Fatalf("NewDecimal128FromString(-.5) string mismatch: %s", got)
	}

	d6, err := NewDecimal128FromString("12345678901234567890e-5")
	if err != nil {
		t.Fatalf("NewDecimal128FromString exponent error: %v", err)
	}
	if got := d6.String(); got != "123456789012345.6789" {
		t.Fatalf("NewDecimal128FromString exponent mismatch: %s", got)
	}

	if _, err := NewDecimal128FromString("bad"); err == nil {
		t.Fatalf("NewDecimal128FromString expected error")
	}
	if _, err := NewDecimal128FromFloat(math.NaN()); err == nil {
		t.Fatalf("NewDecimal128FromFloat NaN expected error")
	}
}

func TestDecimal128Conversions(t *testing.T) {
	d := mustDecimal128(t, "123.0000000000000001")
	intPart, decPart := d.Int64()
	if intPart != 123 || decPart != 1 {
		t.Fatalf("Int64 mismatch: %d %d", intPart, decPart)
	}

	if got := d.String(); got != "123.0000000000000001" {
		t.Fatalf("String mismatch: %s", got)
	}

	d2 := mustDecimal128(t, "1.2")
	if got := d2.StringFixed(4); got != "1.2000" {
		t.Fatalf("StringFixed mismatch: %s", got)
	}
	if got := d2.StringFixed(0); got != "1" {
		t.Fatalf("StringFixed(0) mismatch: %s", got)
	}
	if got16, got20 := d2.StringFixed(16), d2.StringFixed(20); got16 != got20 {
		t.Fatalf("StringFixed n>16 mismatch: %s vs %s", got16, got20)
	}

	if diff := math.Abs(d2.Float64() - 1.2); diff > 1e-9 {
		t.Fatalf("Float64 mismatch: diff=%g", diff)
	}
}

func TestDecimal128Checking(t *testing.T) {
	zero := Decimal128{}
	if !zero.IsZero() {
		t.Fatalf("IsZero failed")
	}
	if zero.IsPositive() || zero.IsNegative() {
		t.Fatalf("zero sign mismatch")
	}
	if zero.Sign() != 0 {
		t.Fatalf("zero Sign mismatch")
	}

	pos := NewDecimal128FromInt(1)
	if !pos.IsPositive() || pos.IsNegative() || pos.Sign() != 1 {
		t.Fatalf("positive sign mismatch")
	}

	neg := NewDecimal128FromInt(-1)
	if !neg.IsNegative() || neg.IsPositive() || neg.Sign() != 2 {
		t.Fatalf("negative sign mismatch")
	}
}

func TestDecimal128Modification(t *testing.T) {
	d := mustDecimal128(t, "1.25")
	if got := d.Neg().String(); got != "-1.25" {
		t.Fatalf("Neg mismatch: %s", got)
	}
	if got := d.Inv().String(); got != "0.8" {
		t.Fatalf("Inv mismatch: %s", got)
	}
	if got := d.Abs().String(); got != "1.25" {
		t.Fatalf("Abs mismatch: %s", got)
	}

	if got := mustDecimal128(t, "123.4567").Truncate(2).String(); got != "123.45" {
		t.Fatalf("Truncate mismatch: %s", got)
	}
	if got := mustDecimal128(t, "123.45").Truncate(-1).String(); got != "120" {
		t.Fatalf("Truncate negative mismatch: %s", got)
	}
	if got := mustDecimal128(t, "1.23").Truncate(17); !got.Equal(mustDecimal128(t, "1.23")) {
		t.Fatalf("Truncate n>16 mismatch: %s", got.String())
	}
	if got := mustDecimal128(t, "1.23").Truncate(-17); !got.IsZero() {
		t.Fatalf("Truncate n<-16 mismatch: %s", got.String())
	}

	if got := mustDecimal128(t, "1.23").Shift(1).String(); got != "12.3" {
		t.Fatalf("Shift(+1) mismatch: %s", got)
	}
	if got := mustDecimal128(t, "1.23").Shift(-1).String(); got != "0.123" {
		t.Fatalf("Shift(-1) mismatch: %s", got)
	}

	if got := mustDecimal128(t, "1.25").Round(1).String(); got != "1.2" {
		t.Fatalf("Round banker mismatch: %s", got)
	}
	if got := mustDecimal128(t, "1.35").Round(1).String(); got != "1.4" {
		t.Fatalf("Round banker mismatch: %s", got)
	}
	if got := mustDecimal128(t, "1.21").RoundAwayFromZero(1).String(); got != "1.3" {
		t.Fatalf("RoundAwayFromZero mismatch: %s", got)
	}
	if got := mustDecimal128(t, "-1.21").RoundAwayFromZero(1).String(); got != "-1.3" {
		t.Fatalf("RoundAwayFromZero negative mismatch: %s", got)
	}
	if got := mustDecimal128(t, "1.29").RoundTowardToZero(1).String(); got != "1.2" {
		t.Fatalf("RoundTowardToZero mismatch: %s", got)
	}
	if got := mustDecimal128(t, "1.2").Ceil(0).String(); got != "2" {
		t.Fatalf("Ceil mismatch: %s", got)
	}
	if got := mustDecimal128(t, "-1.2").Ceil(0).String(); got != "-1" {
		t.Fatalf("Ceil negative mismatch: %s", got)
	}
	if got := mustDecimal128(t, "1.2").Floor(0).String(); got != "1" {
		t.Fatalf("Floor mismatch: %s", got)
	}
	if got := mustDecimal128(t, "-1.2").Floor(0).String(); got != "-2" {
		t.Fatalf("Floor negative mismatch: %s", got)
	}
}

func TestDecimal128Comparison(t *testing.T) {
	a := mustDecimal128(t, "1.5")
	b := mustDecimal128(t, "2")

	if !a.LessThan(b) {
		t.Fatalf("LessThan mismatch")
	}
	if !b.GreaterThan(a) {
		t.Fatalf("GreaterThan mismatch")
	}
	if !a.LessOrEqual(a) {
		t.Fatalf("LessOrEqual mismatch")
	}
	if !b.GreaterOrEqual(b) {
		t.Fatalf("GreaterOrEqual mismatch")
	}
	if !a.Equal(mustDecimal128(t, "1.5")) {
		t.Fatalf("Equal mismatch")
	}
}

func TestDecimal128Arithmetic(t *testing.T) {
	a := mustDecimal128(t, "1.5")
	b := mustDecimal128(t, "2.25")

	if got := a.Add(b).String(); got != "3.75" {
		t.Fatalf("Add mismatch: %s", got)
	}
	if got := b.Sub(a).String(); got != "0.75" {
		t.Fatalf("Sub mismatch: %s", got)
	}
	if got := a.Mul(b).String(); got != "3.375" {
		t.Fatalf("Mul mismatch: %s", got)
	}
	if got := b.Div(a).String(); got != "1.5" {
		t.Fatalf("Div mismatch: %s", got)
	}
	if got := b.Mod(a).String(); got != "0.75" {
		t.Fatalf("Mod mismatch: %s", got)
	}
	if got := a.Div(Decimal128{}).String(); got != "1.5" {
		t.Fatalf("Div by zero mismatch: %s", got)
	}
	if got := a.Mod(Decimal128{}).String(); got != "1.5" {
		t.Fatalf("Mod by zero mismatch: %s", got)
	}
}

func TestDecimal128Transcendental(t *testing.T) {
	if got := NewDecimal128FromInt(2).Pow(NewDecimal128FromInt(3)).String(); got != "8" {
		t.Fatalf("Pow mismatch: %s", got)
	}
	if got := NewDecimal128FromInt(2).Pow(NewDecimal128FromInt(-3)).String(); got != "0.125" {
		t.Fatalf("Pow negative mismatch: %s", got)
	}
	if got := NewDecimal128FromInt(4).Sqrt().String(); got != "2" {
		t.Fatalf("Sqrt mismatch: %s", got)
	}
	if got := NewDecimal128FromInt(-4).Sqrt().String(); got != "-4" {
		t.Fatalf("Sqrt negative mismatch: %s", got)
	}
	if got := (Decimal128{}).Exp().String(); got != "1" {
		t.Fatalf("Exp(0) mismatch: %s", got)
	}
	if got := NewDecimal128FromInt(1).Log().String(); got != "0" {
		t.Fatalf("Log(1) mismatch: %s", got)
	}
	if got := NewDecimal128FromInt(1).Log2().String(); got != "0" {
		t.Fatalf("Log2(1) mismatch: %s", got)
	}
	if got := NewDecimal128FromInt(1).Log10().String(); got != "0" {
		t.Fatalf("Log10(1) mismatch: %s", got)
	}
}

func TestDecimal128EncodeDecode(t *testing.T) {
	d := mustDecimal128(t, "123.456")
	bin, err := d.EncodeBinary()
	if err != nil {
		t.Fatalf("EncodeBinary error: %v", err)
	}
	if len(bin) != 16 {
		t.Fatalf("EncodeBinary length mismatch: %d", len(bin))
	}
	d2, err := NewDecimal128FromBinary(bin)
	if err != nil {
		t.Fatalf("NewDecimal128FromBinary error: %v", err)
	}
	if !d2.Equal(d) {
		t.Fatalf("NewDecimal128FromBinary mismatch: %s", d2.String())
	}
	if _, err := NewDecimal128FromBinary([]byte{1, 2, 3}); err == nil {
		t.Fatalf("NewDecimal128FromBinary expected error")
	}

	jsonBytes, err := d.EncodeJSON()
	if err != nil {
		t.Fatalf("EncodeJSON error: %v", err)
	}
	d3, err := NewDecimal128FromJSON(jsonBytes)
	if err != nil {
		t.Fatalf("NewDecimal128FromJSON string error: %v", err)
	}
	if !d3.Equal(d) {
		t.Fatalf("NewDecimal128FromJSON string mismatch: %s", d3.String())
	}
	d4, err := NewDecimal128FromJSON([]byte("123.456"))
	if err != nil {
		t.Fatalf("NewDecimal128FromJSON number error: %v", err)
	}
	if got := d4.String(); got != "123.456" {
		t.Fatalf("NewDecimal128FromJSON number mismatch: %s", got)
	}
	if _, err := NewDecimal128FromJSON([]byte("\"bad\"")); err == nil {
		t.Fatalf("NewDecimal128FromJSON invalid expected error")
	}
}

func TestDecimal128Append(t *testing.T) {
	d := mustDecimal128(t, "123.456")

	bin, err := d.EncodeBinary()
	if err != nil {
		t.Fatalf("EncodeBinary error: %v", err)
	}
	{
		dst := d.AppendBinary(nil)
		if !bytes.Equal(dst, bin) {
			t.Fatalf("AppendBinary mismatch")
		}
	}

	jsonBytes, err := d.EncodeJSON()
	if err != nil {
		t.Fatalf("EncodeJSON error: %v", err)
	}
	{
		dst := d.AppendJSON(nil)
		if !bytes.Equal(dst, jsonBytes) {
			t.Fatalf("AppendJSON mismatch")
		}
	}

	{
		dst := d.AppendString(nil)
		if got := string(dst); got != d.String() {
			t.Fatalf("AppendString mismatch: %s", got)
		}
	}

	{
		dst := d.AppendStringFixed(nil, 4)
		if got := string(dst); got != d.StringFixed(4) {
			t.Fatalf("AppendStringFixed mismatch: %s", got)
		}
	}
}
