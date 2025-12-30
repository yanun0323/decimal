package decimal

import (
	"bytes"
	"math"
	"testing"
)

func mustDecimal(t *testing.T, s string) Decimal256 {
	t.Helper()
	d, err := New256FromString(s)
	if err != nil {
		t.Fatalf("New256FromString(%q) error: %v", s, err)
	}
	return d
}

func TestDecimal256Constructors(t *testing.T) {
	d1 := New256(123456789, 987654321)
	if got := d1.String(); got != "123456789.987654321" {
		t.Fatalf("New256 string mismatch: %s", got)
	}

	d2 := New256FromInt(-123456789)
	if got := d2.String(); got != "-123456789" {
		t.Fatalf("New256FromInt string mismatch: %s", got)
	}

	d3, err := New256FromFloat(1.25)
	if err != nil {
		t.Fatalf("New256FromFloat error: %v", err)
	}
	if diff := math.Abs(d3.Float64() - 1.25); diff > 1e-9 {
		t.Fatalf("New256FromFloat value mismatch: diff=%g", diff)
	}

	d4, err := New256FromString("  +1_234.4500e0 ")
	if err != nil {
		t.Fatalf("New256FromString error: %v", err)
	}
	if got := d4.String(); got != "1234.45" {
		t.Fatalf("New256FromString string mismatch: %s", got)
	}

	d5, err := New256FromString("-.5")
	if err != nil {
		t.Fatalf("New256FromString(-.5) error: %v", err)
	}
	if got := d5.String(); got != "-0.5" {
		t.Fatalf("New256FromString(-.5) string mismatch: %s", got)
	}

	d6, err := New256FromString("12345678901234567890123456789012345")
	if err != nil {
		t.Fatalf("New256FromString precision error: %v", err)
	}
	if got := d6.String(); got != "45678901234567890123456789012345" {
		t.Fatalf("New256 precision mismatch: %s", got)
	}

	if _, err := New256FromString("bad"); err == nil {
		t.Fatalf("New256FromString expected error")
	}
	if _, err := New256FromFloat(math.NaN()); err == nil {
		t.Fatalf("New256FromFloat NaN expected error")
	}
}

func TestDecimal256Conversions(t *testing.T) {
	d := mustDecimal(t, "123.000000000000000001")
	intPart, decPart := d.Int64()
	if intPart != 123 || decPart != 100000000000000 {
		t.Fatalf("Int64 mismatch: %d %d", intPart, decPart)
	}

	if got := d.String(); got != "123.000000000000000001" {
		t.Fatalf("String mismatch: %s", got)
	}

	d2 := mustDecimal(t, "1.2")
	if got := d2.StringFixed(4); got != "1.2000" {
		t.Fatalf("StringFixed mismatch: %s", got)
	}
	if got := d2.StringFixed(0); got != "1" {
		t.Fatalf("StringFixed(0) mismatch: %s", got)
	}
	if got32, got40 := d2.StringFixed(32), d2.StringFixed(40); got32 != got40 {
		t.Fatalf("StringFixed n>32 mismatch: %s vs %s", got32, got40)
	}

	if diff := math.Abs(d2.Float64() - 1.2); diff > 1e-9 {
		t.Fatalf("Float64 mismatch: diff=%g", diff)
	}
}

func TestDecimal256Checking(t *testing.T) {
	zero := Decimal256{}
	if !zero.IsZero() {
		t.Fatalf("IsZero failed")
	}
	if zero.IsPositive() || zero.IsNegative() {
		t.Fatalf("zero sign mismatch")
	}
	if zero.Sign() != 0 {
		t.Fatalf("zero Sign mismatch")
	}

	pos := New256FromInt(1)
	if !pos.IsPositive() || pos.IsNegative() || pos.Sign() != 1 {
		t.Fatalf("positive sign mismatch")
	}

	neg := New256FromInt(-1)
	if !neg.IsNegative() || neg.IsPositive() || neg.Sign() != 2 {
		t.Fatalf("negative sign mismatch")
	}
}

func TestDecimal256Modification(t *testing.T) {
	d := mustDecimal(t, "1.25")
	if got := d.Neg().String(); got != "-1.25" {
		t.Fatalf("Neg mismatch: %s", got)
	}
	if got := d.Inv().String(); got != "0.8" {
		t.Fatalf("Inv mismatch: %s", got)
	}
	if got := d.Abs().String(); got != "1.25" {
		t.Fatalf("Abs mismatch: %s", got)
	}

	if got := mustDecimal(t, "123.4567").Truncate(2).String(); got != "123.45" {
		t.Fatalf("Truncate mismatch: %s", got)
	}
	if got := mustDecimal(t, "123.45").Truncate(-1).String(); got != "120" {
		t.Fatalf("Truncate negative mismatch: %s", got)
	}
	if got := mustDecimal(t, "1.23").Truncate(33); !got.Equal(mustDecimal(t, "1.23")) {
		t.Fatalf("Truncate n>32 mismatch: %s", got.String())
	}
	if got := mustDecimal(t, "1.23").Truncate(-33); !got.IsZero() {
		t.Fatalf("Truncate n<-32 mismatch: %s", got.String())
	}

	if got := mustDecimal(t, "1.23").Shift(1).String(); got != "12.3" {
		t.Fatalf("Shift(+1) mismatch: %s", got)
	}
	if got := mustDecimal(t, "1.23").Shift(-1).String(); got != "0.123" {
		t.Fatalf("Shift(-1) mismatch: %s", got)
	}

	if got := mustDecimal(t, "1.25").Round(1).String(); got != "1.2" {
		t.Fatalf("Round banker mismatch: %s", got)
	}
	if got := mustDecimal(t, "1.35").Round(1).String(); got != "1.4" {
		t.Fatalf("Round banker mismatch: %s", got)
	}
	if got := mustDecimal(t, "1.21").RoundAwayFromZero(1).String(); got != "1.3" {
		t.Fatalf("RoundAwayFromZero mismatch: %s", got)
	}
	if got := mustDecimal(t, "-1.21").RoundAwayFromZero(1).String(); got != "-1.3" {
		t.Fatalf("RoundAwayFromZero negative mismatch: %s", got)
	}
	if got := mustDecimal(t, "1.29").RoundTowardToZero(1).String(); got != "1.2" {
		t.Fatalf("RoundTowardToZero mismatch: %s", got)
	}
	if got := mustDecimal(t, "1.2").Ceil(0).String(); got != "2" {
		t.Fatalf("Ceil mismatch: %s", got)
	}
	if got := mustDecimal(t, "-1.2").Ceil(0).String(); got != "-1" {
		t.Fatalf("Ceil negative mismatch: %s", got)
	}
	if got := mustDecimal(t, "1.2").Floor(0).String(); got != "1" {
		t.Fatalf("Floor mismatch: %s", got)
	}
	if got := mustDecimal(t, "-1.2").Floor(0).String(); got != "-2" {
		t.Fatalf("Floor negative mismatch: %s", got)
	}
}

func TestDecimal256Comparison(t *testing.T) {
	a := mustDecimal(t, "1.5")
	b := mustDecimal(t, "2")

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
	if !a.Equal(mustDecimal(t, "1.5")) {
		t.Fatalf("Equal mismatch")
	}
}

func TestDecimal256Arithmetic(t *testing.T) {
	a := mustDecimal(t, "1.5")
	b := mustDecimal(t, "2.25")

	if got := a.Add(b).String(); got != "3.75" {
		t.Fatalf("Add mismatch: %s", got)
	}
	if got := b.Sub(a).String(); got != "0.75" {
		t.Fatalf("Sub mismatch: %s", got)
	}
	if got := a.Mul(New256FromInt(2)).String(); got != "3" {
		t.Fatalf("Mul mismatch: %s", got)
	}
	if got := mustDecimal(t, "3").Div(New256FromInt(2)).String(); got != "1.5" {
		t.Fatalf("Div mismatch: %s", got)
	}
	if got := mustDecimal(t, "5.5").Mod(New256FromInt(2)).String(); got != "1.5" {
		t.Fatalf("Mod mismatch: %s", got)
	}
	if got := a.Div(Decimal256{}).String(); got != "1.5" {
		t.Fatalf("Div by zero mismatch: %s", got)
	}
	if got := a.Mod(Decimal256{}).String(); got != "1.5" {
		t.Fatalf("Mod by zero mismatch: %s", got)
	}
}

func TestDecimal256Transcendental(t *testing.T) {
	if got := New256FromInt(2).Pow(New256FromInt(3)).String(); got != "8" {
		t.Fatalf("Pow mismatch: %s", got)
	}
	if got := New256FromInt(2).Pow(New256FromInt(-3)).String(); got != "0.125" {
		t.Fatalf("Pow negative mismatch: %s", got)
	}
	if got := New256FromInt(4).Sqrt().String(); got != "2" {
		t.Fatalf("Sqrt mismatch: %s", got)
	}
	if got := New256FromInt(-4).Sqrt().String(); got != "-4" {
		t.Fatalf("Sqrt negative mismatch: %s", got)
	}
	if got := (Decimal256{}).Exp().String(); got != "1" {
		t.Fatalf("Exp(0) mismatch: %s", got)
	}
	if got := New256FromInt(1).Log().String(); got != "0" {
		t.Fatalf("Log(1) mismatch: %s", got)
	}
	if got := New256FromInt(1).Log2().String(); got != "0" {
		t.Fatalf("Log2(1) mismatch: %s", got)
	}
	if got := New256FromInt(1).Log10().String(); got != "0" {
		t.Fatalf("Log10(1) mismatch: %s", got)
	}
}

func TestDecimal256EncodeDecode(t *testing.T) {
	d := mustDecimal(t, "123.456")
	bin, err := d.EncodeBinary()
	if err != nil {
		t.Fatalf("EncodeBinary error: %v", err)
	}
	if len(bin) != 32 {
		t.Fatalf("EncodeBinary length mismatch: %d", len(bin))
	}
	d2, err := New256FromBinary(bin)
	if err != nil {
		t.Fatalf("New256FromBinary error: %v", err)
	}
	if !d2.Equal(d) {
		t.Fatalf("New256FromBinary mismatch: %s", d2.String())
	}
	if _, err := New256FromBinary([]byte{1, 2, 3}); err == nil {
		t.Fatalf("New256FromBinary expected error")
	}

	jsonBytes, err := d.EncodeJSON()
	if err != nil {
		t.Fatalf("EncodeJSON error: %v", err)
	}
	d3, err := New256FromJSON(jsonBytes)
	if err != nil {
		t.Fatalf("New256FromJSON string error: %v", err)
	}
	if !d3.Equal(d) {
		t.Fatalf("New256FromJSON string mismatch: %s", d3.String())
	}
	d4, err := New256FromJSON([]byte("123.456"))
	if err != nil {
		t.Fatalf("New256FromJSON number error: %v", err)
	}
	if got := d4.String(); got != "123.456" {
		t.Fatalf("New256FromJSON number mismatch: %s", got)
	}
	if _, err := New256FromJSON([]byte("\"bad\"")); err == nil {
		t.Fatalf("New256FromJSON invalid expected error")
	}
}

func TestDecimal256Append(t *testing.T) {
	d := mustDecimal(t, "123.456")

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
