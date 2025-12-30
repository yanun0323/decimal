package decimal

// Don't use this. just for definition.
type Decimal interface {
	// conversions
	Int64() (intPart, decimalPart int64)
	Float64() float64
	String() string
	StringFixed(int) string

	// checking
	IsZero() bool
	IsPositive() bool
	IsNegative() bool
	Sign() int // 0:zero 1:pos 2:neg

	// modification
	Neg() Decimal
	Inv() Decimal
	Abs() Decimal
	Truncate(int) Decimal
	Shift(int) Decimal
	Round(int) Decimal
	RoundAwayFromZero(int) Decimal
	RoundTowardToZero(int) Decimal
	Ceil(int) Decimal
	Floor(int) Decimal

	// comparison
	Equal(Decimal) bool
	GreaterThan(Decimal) bool
	LessThan(Decimal) bool
	GreaterOrEqual(Decimal) bool
	LessOrEqual(Decimal) bool

	// arithmetic operations
	Add(Decimal) Decimal
	Sub(Decimal) Decimal
	Mul(Decimal) Decimal
	Div(Decimal) Decimal
	Mod(Decimal) Decimal

	// transcendental operations
	Pow(Decimal) Decimal
	Sqrt() Decimal
	Exp() Decimal
	Log() Decimal
	Log2() Decimal
	Log10() Decimal

	// normalized
	EncodeBinary() ([]byte, error)
	EncodeJSON() ([]byte, error)
	AppendBinary([]byte) []byte
	AppendJSON([]byte) []byte
	AppendString([]byte) []byte
	AppendStringFixed([]byte, int) []byte
}
