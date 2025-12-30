# Decimal256

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)

A zero-allocation, fixed-scale decimal library built on a 256-bit two's-complement integer.

## Requirements

- Go 1.25+ (per `go.mod`)

## Import

```go
import "github.com/yanun0323/decimal"
```

## Overview

Decimal256 is a fixed-scale decimal type with **32 fractional digits**.

- **Value model**: `raw / 10^32`
- **Storage**: 256-bit two's-complement integer (4 x uint64)
- **Overflow**: truncated to 256-bit (wrap-around)
- **Zero value**: valid and represents `0`
- **No big.Int**: all arithmetic is fixed-size

## Constructors

- `NewDecimal256(intPart, decimalPart int64) Decimal256`
  - `decimalPart` is treated as **fractional digits**.
  - Example: `NewDecimal256(123, 45)` = `123.45`.
  - More than 32 digits are truncated toward zero.
- `NewDecimal256FromString(string) (Decimal256, error)`
  - Accepts leading/trailing ASCII whitespace, `_` separators, optional `.` and `e/E`.
  - Excess fractional digits are truncated to 32 places.
- `NewDecimal256FromInt(int64) Decimal256`
- `NewDecimal256FromFloat(float64) (Decimal256, error)`
  - Truncates toward zero. `NaN`/`Inf` return error.
- `NewDecimal256FromBinary([]byte) (Decimal256, error)`
  - Expects 32 bytes, little-endian.
- `NewDecimal256FromJSON([]byte) (Decimal256, error)`
  - Accepts JSON **string** or **number**.

## Conversions & Formatting

- `Int64() (intPart, decimalPart int64)`
  - `decimalPart` is scaled by 10^32 (same internal scale).
- `Float64() float64`
- `String() string`
  - Removes trailing zeros in the fractional part.
- `StringFixed(n int) string`
  - If `n > 32`, it is truncated to 32. If `n <= 0`, only the integer part is returned.

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

Rules for digit operations:

- If `n > 32`: **no change**
- If `n <= -32`: **return zero**
- If `n < 0`: operation applies to integer digits (e.g. `Truncate(-1)` => tens)

## Transcendental

- `Pow` (integer exponent only; exponent is truncated toward zero)
- `Sqrt`, `Exp`, `Log`, `Log2`, `Log10`
  - For invalid input (e.g. negative for `Sqrt`, non-positive for `Log`), the original value is returned.

## Binary / JSON Encoding

- **Binary**: 32 bytes, little-endian (`[4]uint64` word order)
- **JSON**: encoded as string; decoder accepts string or number

## Performance

Designed for zero allocations and predictable cost. For benchmarks:

```sh
go test -bench . ./...
```

## Contributing

Pull requests are welcome. For major changes, open an issue to discuss the design first.
