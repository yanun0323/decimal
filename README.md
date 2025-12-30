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

## Predefined values

- `Zero128`, `One128`, `Ten128`, `Hundred128`
- `Zero256`, `One256`, `Ten256`, `Hundred256`
- `Zero512`, `One512`, `Ten512`, `Hundred512`

## Constructors (per type)

Replace `XXX` with `128`, `256`, or `512`:

- `NewXXX(intPart, decimalPart int64) DecimalXXX`
  - `decimalPart` is treated as fractional digits.
  - Example: `New256(123, 45)` = `123.45`.
  - Excess fractional digits are truncated toward zero.
  - Precision rules apply (integer low *N*, fractional high *N*).
- `NewXXXFromString(string) (DecimalXXX, error)`
  - Accepts leading/trailing ASCII whitespace, `_` separators, optional `.` and `e/E`.
  - Exponent shifting is applied first, then precision rules apply.
- `NewXXXFromInt(int64) DecimalXXX`
  - Precision rules apply (integer low *N*).
- `NewXXXFromFloat(float64) (DecimalXXX, error)`
  - Truncates toward zero. `NaN`/`Inf` return error.
  - Precision rules apply after conversion.
- `NewXXXFromBinary([]byte) (DecimalXXX, error)`
  - Expects fixed length, little-endian (see Binary section).
  - Precision rules apply after decoding.
- `NewXXXFromJSON([]byte) (DecimalXXX, error)`
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

## Database integration

- SQL (`database/sql`)
  - Implements `sql.Scanner` and `driver.Valuer`
  - `Scan` accepts `string`, `[]byte`, `int64`, `float64`, and `NULL` (NULL becomes zero value)
  - `Value` returns the decimal string
- MongoDB Go Driver v2 (`go.mongodb.org/mongo-driver/v2`)
  - Implements `bson.ValueMarshaler` / `bson.ValueUnmarshaler`
  - `Decimal128` is encoded as BSON Decimal128
  - `Decimal256` / `Decimal512` are encoded as BSON string

## Generic interface

The package also exposes a generic interface for compile-time constraints:

- `type Decimal[T decimal] interface { ... }`

## Performance

Designed for zero allocations and predictable cost. For benchmarks:

```markdown
go test -bench=. -run=none -benchmem --count=1 ./...
goos: darwin
goarch: arm64
pkg: github.com/yanun0323/decimal
cpu: Apple M2
Decimal128/New128               225.0 ns/op    0 B/op   0 allocs/op
Decimal128/New128FromString     403.4 ns/op    0 B/op   0 allocs/op
Decimal128/New128FromInt        199.9 ns/op    0 B/op   0 allocs/op
Decimal128/New128FromFloat      207.8 ns/op    0 B/op   0 allocs/op
Decimal128/Int64                157.7 ns/op    0 B/op   0 allocs/op
Decimal128/Float64              14.36 ns/op    0 B/op   0 allocs/op
Decimal128/String               261.2 ns/op   40 B/op   3 allocs/op
Decimal128/StringFixed          276.4 ns/op   48 B/op   3 allocs/op
Decimal128/IsZero               1.953 ns/op    0 B/op   0 allocs/op
Decimal128/IsPositive           2.278 ns/op    0 B/op   0 allocs/op
Decimal128/IsNegative           2.229 ns/op    0 B/op   0 allocs/op
Decimal128/Sign                 2.301 ns/op    0 B/op   0 allocs/op
Decimal128/Neg                  2.190 ns/op    0 B/op   0 allocs/op
Decimal128/Inv                  318.3 ns/op    0 B/op   0 allocs/op
Decimal128/Abs                  4.024 ns/op    0 B/op   0 allocs/op
Decimal128/Truncate             441.5 ns/op    0 B/op   0 allocs/op
Decimal128/Shift                10.06 ns/op    0 B/op   0 allocs/op
Decimal128/Round                452.0 ns/op    0 B/op   0 allocs/op
Decimal128/RoundAwayFromZero    450.2 ns/op    0 B/op   0 allocs/op
Decimal128/RoundTowardToZero    444.8 ns/op    0 B/op   0 allocs/op
Decimal128/Ceil                 254.7 ns/op    0 B/op   0 allocs/op
Decimal128/Floor                247.1 ns/op    0 B/op   0 allocs/op
Decimal128/Equal                1.958 ns/op    0 B/op   0 allocs/op
Decimal128/GreaterThan          3.386 ns/op    0 B/op   0 allocs/op
Decimal128/LessThan             3.325 ns/op    0 B/op   0 allocs/op
Decimal128/GreaterOrEqual       4.156 ns/op    0 B/op   0 allocs/op
Decimal128/LessOrEqual          4.020 ns/op    0 B/op   0 allocs/op
Decimal128/Add                  3.041 ns/op    0 B/op   0 allocs/op
Decimal128/Sub                  2.988 ns/op    0 B/op   0 allocs/op
Decimal128/Mul                  748.2 ns/op    0 B/op   0 allocs/op
Decimal128/Div                  569.0 ns/op    0 B/op   0 allocs/op
Decimal128/Mod                  149.4 ns/op    0 B/op   0 allocs/op
Decimal128/Pow                   2691 ns/op    0 B/op   0 allocs/op
Decimal128/Sqrt                  1355 ns/op    0 B/op   0 allocs/op
Decimal128/Exp                   6213 ns/op    0 B/op   0 allocs/op
Decimal128/Log                   3389 ns/op    0 B/op   0 allocs/op
Decimal128/Log2                  4131 ns/op    0 B/op   0 allocs/op
Decimal128/Log10                 4118 ns/op    0 B/op   0 allocs/op
Decimal128/EncodeBinary         10.89 ns/op   16 B/op   1 allocs/op
Decimal128/AppendBinary         2.657 ns/op    0 B/op   0 allocs/op
Decimal128/NewFromBinary        186.3 ns/op    0 B/op   0 allocs/op
Decimal128/EncodeJSON           356.9 ns/op   64 B/op   4 allocs/op
Decimal128/AppendJSON           226.1 ns/op    0 B/op   0 allocs/op
Decimal128/NewFromJSON          385.5 ns/op    0 B/op   0 allocs/op
Decimal128/AppendString         217.4 ns/op    0 B/op   0 allocs/op
Decimal128/AppendStringFixed    219.2 ns/op    0 B/op   0 allocs/op

Decimal256/New256               598.6 ns/op    0 B/op   0 allocs/op
Decimal256/New256FromString     822.7 ns/op    0 B/op   0 allocs/op
Decimal256/New256FromInt        552.5 ns/op    0 B/op   0 allocs/op
Decimal256/New256FromFloat      323.7 ns/op    0 B/op   0 allocs/op
Decimal256/Int64                452.3 ns/op    0 B/op   0 allocs/op
Decimal256/Float64              51.36 ns/op    0 B/op   0 allocs/op
Decimal256/String               633.9 ns/op   72 B/op   3 allocs/op
Decimal256/StringFixed          618.7 ns/op   80 B/op   3 allocs/op
Decimal256/IsZero               2.222 ns/op    0 B/op   0 allocs/op
Decimal256/IsPositive           3.396 ns/op    0 B/op   0 allocs/op
Decimal256/IsNegative           3.373 ns/op    0 B/op   0 allocs/op
Decimal256/Sign                 3.430 ns/op    0 B/op   0 allocs/op
Decimal256/Neg                  6.095 ns/op    0 B/op   0 allocs/op
Decimal256/Inv                   1186 ns/op    0 B/op   0 allocs/op
Decimal256/Abs                  7.581 ns/op    0 B/op   0 allocs/op
Decimal256/Truncate             950.3 ns/op    0 B/op   0 allocs/op
Decimal256/Shift                33.35 ns/op    0 B/op   0 allocs/op
Decimal256/Round                950.5 ns/op    0 B/op   0 allocs/op
Decimal256/RoundAwayFromZero    962.0 ns/op    0 B/op   0 allocs/op
Decimal256/RoundTowardToZero    952.2 ns/op    0 B/op   0 allocs/op
Decimal256/Ceil                 621.6 ns/op    0 B/op   0 allocs/op
Decimal256/Floor                612.2 ns/op    0 B/op   0 allocs/op
Decimal256/Equal                1.958 ns/op    0 B/op   0 allocs/op
Decimal256/GreaterThan          5.593 ns/op    0 B/op   0 allocs/op
Decimal256/LessThan             5.599 ns/op    0 B/op   0 allocs/op
Decimal256/GreaterOrEqual       5.949 ns/op    0 B/op   0 allocs/op
Decimal256/LessOrEqual          6.048 ns/op    0 B/op   0 allocs/op
Decimal256/Add                  5.066 ns/op    0 B/op   0 allocs/op
Decimal256/Sub                  5.050 ns/op    0 B/op   0 allocs/op
Decimal256/Mul                   1987 ns/op    0 B/op   0 allocs/op
Decimal256/Div                   1970 ns/op    0 B/op   0 allocs/op
Decimal256/Mod                  432.4 ns/op    0 B/op   0 allocs/op
Decimal256/Pow                  10161 ns/op    0 B/op   0 allocs/op
Decimal256/Sqrt                  6263 ns/op    0 B/op   0 allocs/op
Decimal256/Exp                  29479 ns/op    0 B/op   0 allocs/op
Decimal256/Log                  18032 ns/op    0 B/op   0 allocs/op
Decimal256/Log2                 19972 ns/op    0 B/op   0 allocs/op
Decimal256/Log10                19997 ns/op    0 B/op   0 allocs/op
Decimal256/EncodeBinary         13.18 ns/op   32 B/op   1 allocs/op
Decimal256/AppendBinary         3.421 ns/op    0 B/op   0 allocs/op
Decimal256/NewFromBinary        523.1 ns/op    0 B/op   0 allocs/op
Decimal256/EncodeJSON           726.2 ns/op   96 B/op   4 allocs/op
Decimal256/AppendJSON           577.5 ns/op    0 B/op   0 allocs/op
Decimal256/NewFromJSON          787.0 ns/op    0 B/op   0 allocs/op
Decimal256/AppendString         568.4 ns/op    0 B/op   0 allocs/op
Decimal256/AppendStringFixed    573.8 ns/op    0 B/op   0 allocs/op

Decimal512/New512                1409 ns/op    0 B/op   0 allocs/op
Decimal512/New512FromString      1686 ns/op    0 B/op   0 allocs/op
Decimal512/New512FromInt         1287 ns/op    0 B/op   0 allocs/op
Decimal512/New512FromFloat      730.9 ns/op    0 B/op   0 allocs/op
Decimal512/Int64                 1002 ns/op    0 B/op   0 allocs/op
Decimal512/Float64              47.74 ns/op    0 B/op   0 allocs/op
Decimal512/String                1281 ns/op  104 B/op   3 allocs/op
Decimal512/StringFixed           1260 ns/op  112 B/op   3 allocs/op
Decimal512/IsZero               3.356 ns/op    0 B/op   0 allocs/op
Decimal512/IsPositive           5.852 ns/op    0 B/op   0 allocs/op
Decimal512/IsNegative           5.899 ns/op    0 B/op   0 allocs/op
Decimal512/Sign                 5.889 ns/op    0 B/op   0 allocs/op
Decimal512/Neg                  13.52 ns/op    0 B/op   0 allocs/op
Decimal512/Inv                   6009 ns/op    0 B/op   0 allocs/op
Decimal512/Abs                  15.78 ns/op    0 B/op   0 allocs/op
Decimal512/Truncate              2071 ns/op    0 B/op   0 allocs/op
Decimal512/Shift                111.5 ns/op    0 B/op   0 allocs/op
Decimal512/Round                 2079 ns/op    0 B/op   0 allocs/op
Decimal512/RoundAwayFromZero     2065 ns/op    0 B/op   0 allocs/op
Decimal512/RoundTowardToZero     2051 ns/op    0 B/op   0 allocs/op
Decimal512/Ceil                  1409 ns/op    0 B/op   0 allocs/op
Decimal512/Floor                 1392 ns/op    0 B/op   0 allocs/op
Decimal512/Equal                4.160 ns/op    0 B/op   0 allocs/op
Decimal512/GreaterThan          7.996 ns/op    0 B/op   0 allocs/op
Decimal512/LessThan             7.971 ns/op    0 B/op   0 allocs/op
Decimal512/GreaterOrEqual       10.72 ns/op    0 B/op   0 allocs/op
Decimal512/LessOrEqual          10.76 ns/op    0 B/op   0 allocs/op
Decimal512/Add                  17.12 ns/op    0 B/op   0 allocs/op
Decimal512/Sub                  17.20 ns/op    0 B/op   0 allocs/op
Decimal512/Mul                   8315 ns/op    0 B/op   0 allocs/op
Decimal512/Div                   7221 ns/op    0 B/op   0 allocs/op
Decimal512/Mod                  962.7 ns/op    0 B/op   0 allocs/op
Decimal512/Pow                  27151 ns/op    0 B/op   0 allocs/op
Decimal512/Sqrt                 30268 ns/op    0 B/op   0 allocs/op
Decimal512/Exp                 168830 ns/op    0 B/op   0 allocs/op
Decimal512/Log                 110280 ns/op    0 B/op   0 allocs/op
Decimal512/Log2                117589 ns/op    0 B/op   0 allocs/op
Decimal512/Log10               117869 ns/op    0 B/op   0 allocs/op
Decimal512/EncodeBinary         21.48 ns/op   64 B/op   1 allocs/op
Decimal512/AppendBinary         16.96 ns/op    0 B/op   0 allocs/op
Decimal512/NewFromBinary         1190 ns/op    0 B/op   0 allocs/op
Decimal512/EncodeJSON            1404 ns/op  128 B/op   4 allocs/op
Decimal512/AppendJSON            1255 ns/op    0 B/op   0 allocs/op
Decimal512/NewFromJSON           1686 ns/op    0 B/op   0 allocs/op
Decimal512/AppendString          1242 ns/op    0 B/op   0 allocs/op
Decimal512/AppendStringFixed     1226 ns/op    0 B/op   0 allocs/op
```

## Contributing

Pull requests are welcome. For major changes, open an issue to discuss the design first.
