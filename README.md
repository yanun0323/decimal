# Decimal

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)

A zero-allocation, fixed-scale decimal library with three fixed-size implementations. `Decimal128` / `Decimal256` / `Decimal512`

## Requirements

- Go 1.25+ (per `go.mod`)

## Import

```go
import "github.com/yanun0323/decimal"
```

## Overview

Each type has its own fixed scale and precision:

- `Decimal128`: scale `10^16`, integer **16** digits, fractional **16** digits
- `Decimal256`: scale `10^32`, integer **32** digits, fractional **32** digits
- `Decimal512`: scale `10^64`, integer **64** digits, fractional **64** digits

Common properties:

- **Value model**: `raw / 10^scale`
- **Overflow**: truncated to the underlying bit width (wrap-around)
- **Zero value**: valid and represents `0`
- **No big.Int**: all arithmetic is fixed-size

Memory layout:

- `Decimal128`: 128-bit two's-complement integer (2 x uint64)
- `Decimal256`: 256-bit two's-complement integer (4 x uint64)
- `Decimal512`: 512-bit two's-complement integer (8 x uint64)

### Precision rules

For all constructors/parsers:

- **Integer part**: keep only the lowest *N* digits (higher digits are dropped)
- **Fractional part**: keep only the highest *N* digits (lower digits are dropped)

Where *N* is the fractional precision for the type (16/32/64).

For string/JSON parsing, **exponent shifting is applied first**, then the precision rules are applied.

## Constructors (per type)

Replace `XXX` with `128`, `256`, or `512`:

- `NewDecimalXXX(intPart, decimalPart int64) DecimalXXX`
  - `decimalPart` is treated as fractional digits.
  - Example: `NewDecimal256(123, 45)` = `123.45`.
  - Excess fractional digits are truncated toward zero.
  - Precision rules apply (integer low *N*, fractional high *N*).
- `NewDecimalXXXFromString(string) (DecimalXXX, error)`
  - Accepts leading/trailing ASCII whitespace, `_` separators, optional `.` and `e/E`.
  - Exponent shifting is applied first, then precision rules apply.
- `NewDecimalXXXFromInt(int64) DecimalXXX`
  - Precision rules apply (integer low *N*).
- `NewDecimalXXXFromFloat(float64) (DecimalXXX, error)`
  - Truncates toward zero. `NaN`/`Inf` return error.
  - Precision rules apply after conversion.
- `NewDecimalXXXFromBinary([]byte) (DecimalXXX, error)`
  - Expects fixed length, little-endian (see Binary section).
  - Precision rules apply after decoding.
- `NewDecimalXXXFromJSON([]byte) (DecimalXXX, error)`
  - Accepts JSON **string** or **number**.
  - Exponent shifting is applied first, then precision rules apply.

## Conversions & Formatting

- `Int64() (intPart, decimalPart int64)`
  - `decimalPart` is scaled by the type scale (`10^16`, `10^32`, `10^64`).
- `Float64() float64`
- `String() string`
  - Removes trailing zeros in the fractional part.
- `StringFixed(n int) string`
  - If `n > scaleDigits`, it is truncated to the type scale (16/32/64).
  - If `n <= 0`, only the integer part is returned.

### Zero-allocation append

- `AppendBinary(dst []byte) []byte`
- `AppendJSON(dst []byte) []byte`
- `AppendString(dst []byte) []byte`
- `AppendStringFixed(dst []byte, n int) []byte`

These append into a caller-provided buffer so allocation is fully controlled by the caller.

## Checks

- `IsZero()`, `IsPositive()`, `IsNegative()`, `Sign()`

## Arithmetic & Comparison

- `Add`, `Sub`, `Mul`, `Div`, `Mod`
  - `Div`/`Mod` by zero return the original value.
- `Equal`, `GreaterThan`, `LessThan`, `GreaterOrEqual`, `LessOrEqual`

## Rounding & Modification

- `Neg`, `Inv`, `Abs`, `Truncate`, `Shift`
- `Round` (banker's rounding)
- `RoundAwayFromZero`, `RoundTowardToZero`
- `Ceil`, `Floor`

Rules for digit operations (scaleDigits = 16/32/64 by type):

- If `n > scaleDigits`: **no change**
- If `n <= -scaleDigits`: **return zero**
- If `n < 0`: operation applies to integer digits (e.g. `Truncate(-1)` => tens)

## Transcendental

- `Pow` (integer exponent only; exponent is truncated toward zero)
- `Sqrt`, `Exp`, `Log`, `Log2`, `Log10`
  - For invalid input (e.g. negative for `Sqrt`, non-positive for `Log`), the original value is returned.

## Binary / JSON Encoding

- **Binary**: little-endian fixed size
  - `Decimal128`: 16 bytes
  - `Decimal256`: 32 bytes
  - `Decimal512`: 64 bytes
- **JSON**: encoded as string; decoder accepts string or number

## Generic interface

The package also exposes a generic interface for compile-time constraints:

- `type Decimal[T decimal] interface { ... }`

## Performance

Designed for zero allocations and predictable cost. For benchmarks:

```sh
go test -bench . ./...
```

## Contributing

Pull requests are welcome. For major changes, open an issue to discuss the design first.
