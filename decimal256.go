// All overflow will be truncated
package decimal

// Decimal256 store a decimal with int part 32 digit and decimal part 32 digit
//
// zero value is ready to use as zero
type Decimal256 [4]byte

// constructors
func NewDecimal256(intPart, decimalPart int64) Decimal256
func NewDecimal256FromString(string) (Decimal256, error)
func NewDecimal256FromInt(int64) Decimal256
func NewDecimal256FromFloat(float64) Decimal256

// conversions
func (Decimal256) Int64() (intPart, decimalPart int64)
func (Decimal256) Float64() float64
func (Decimal256) String() string
func (Decimal256) StringFixed(int) string

// checking
func (Decimal256) IsZero() bool
func (Decimal256) IsPositive() bool
func (Decimal256) IsNegative() bool
func (Decimal256) Sign() int // 0:zero 1:pos 2:neg

// modification
func (Decimal256) Neg() Decimal256
func (Decimal256) Inv() Decimal256
func (Decimal256) Abs() Decimal256
func (Decimal256) Truncate(int) Decimal256
func (Decimal256) Shift(int) Decimal256
func (Decimal256) Round(int) Decimal256
func (Decimal256) RoundAwayFromZero(int) Decimal256
func (Decimal256) RoundTowardToZero(int) Decimal256
func (Decimal256) Ceil(int) Decimal256
func (Decimal256) Floor(int) Decimal256

// comparison
func (Decimal256) Equal(Decimal256) bool
func (Decimal256) GreaterThan(Decimal256) bool
func (Decimal256) LessThan(Decimal256) bool
func (Decimal256) GreaterOrEqual(Decimal256) bool
func (Decimal256) LessOrEqual(Decimal256) bool

// arithmetic operations
func (Decimal256) Add(Decimal256) Decimal256
func (Decimal256) Sub(Decimal256) Decimal256
func (Decimal256) Mul(Decimal256) Decimal256
func (Decimal256) Div(Decimal256) Decimal256
func (Decimal256) Mod(Decimal256) Decimal256

// transcendental operations
func (Decimal256) Pow(Decimal256) Decimal256
func (Decimal256) Sqrt() Decimal256
func (Decimal256) Exp() Decimal256
func (Decimal256) Log() Decimal256
func (Decimal256) Log2() Decimal256
func (Decimal256) Log10() Decimal256
