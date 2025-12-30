package decimal

var _ Decimal[Decimal128] = (*Decimal128)(nil)
var _ Decimal[Decimal256] = (*Decimal256)(nil)
var _ Decimal[Decimal512] = (*Decimal512)(nil)

type decimal interface {
	Decimal128 | Decimal256 | Decimal512
}

// Definition represents an interface for compile-time constraints
type Decimal[T decimal] interface {
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
	Neg() T
	Inv() T
	Abs() T
	Truncate(int) T
	Shift(int) T
	Round(int) T
	RoundAwayFromZero(int) T
	RoundTowardToZero(int) T
	Ceil(int) T
	Floor(int) T

	// comparison
	Equal(T) bool
	GreaterThan(T) bool
	LessThan(T) bool
	GreaterOrEqual(T) bool
	LessOrEqual(T) bool

	// arithmetic operations
	Add(T) T
	Sub(T) T
	Mul(T) T
	Div(T) T
	Mod(T) T

	// transcendental operations
	Pow(T) T
	Sqrt() T
	Exp() T
	Log() T
	Log2() T
	Log10() T

	// normalized
	EncodeBinary() ([]byte, error)
	EncodeJSON() ([]byte, error)
	AppendBinary([]byte) []byte
	AppendJSON([]byte) []byte
	AppendString([]byte) []byte
	AppendStringFixed([]byte, int) []byte
}
