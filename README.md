# Decimal

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)

A super efficient, memory-optimized decimal library based on string type

## Requirements

#### _GO 1.21 or higher_

## Import

```go
import "github.com/yanun0323/decimal"
```

## Features

- The zero-value is 0, and is safe to use without initialization
- Memory-optimized with extremely low memory footprint while maintaining high performance
- Addition, subtraction with no loss of precision
- Database/sql serialization/deserialization
- JSON and XML serialization/deserialization as string
- Fully compatible with [shopspring/decimal](https://github.com/shopspring/decimal) API - all functions are implemented to support the same interface
- Any differences or unimplemented features are documented in the [API Differences](#api-differences) section below

## Supported

- Initialization like new from string, int, int32, float, float64, big.Int
- Addition
- Subtraction
- Multiplication
- Division
- Negative
- Truncate
- Shift
- Compare like equal, greater, less, greater or equal, less or equal
- Round like round, ceil, floor, round bank, round away from zero, round toward to zero

## Usage

```go
// create decimal
zero := decimal.Zero()

d1, err := decimal.New("100,000.555")

d2 := decimal.Require("50_000.05")

// operation
add := d1.Add(d2).String()
println(add)            // 150000.605

sub := d1.Sub(d2).String()
println(sub)            // 50000.505

mul := d1.Mul(d2).String()
println(mul)            // 5000032750.02775

div := d1.Div(d2).String()
println(div)            // 19.9999110009

shift := d1.Shift(-2).String()
println(shift)          // 1000.00555

neg := d1.Neg().String()
println(neg)            // -150000.605

abs := neg.Abs().String()
println(abs)            // 150000.605

truncate := d1.Truncate(1).String()
println(truncate)       // 100000.5

// compare
d1.IsZero()             // false
d1.IsPositive()         // true
d1.IsNegative()         // true

d1.Equal(d2)            // false
d1.Greater(d2)          // true
d1.Less(d2)             // false
d1.GreaterOrEqual(d2)   // true
d1.LessOrEqual(d2)      // false


// method chain
result := d1.Sub(d2).Shift(-5).Add(d1).Truncate(3).String()
```

## Benchmark

Compare to [github.com/shopspring/decimal](https://github.com/shopspring/decimal)

- **Overall Speed**: 1.9-6.5x faster across all operations
- **Memory Efficiency**: 70-88% reduction in memory allocations
- **Creation Operations**: New is 3x faster, NewFromFloat is 6.5x faster with significant memory savings
- **Rounding Operations**: All rounding methods show 3-4x performance improvements with substantial memory reductions

The benchmarks demonstrate consistent performance advantages across creation, arithmetic, transformations, and comparisons while maintaining full API compatibility with shopspring/decimal.

```
New/ShopSpring                 6647031   166.4 ns/op   200 B/op    7 allocs/op
New/Decimal                   23945144   55.54 ns/op    24 B/op    1 allocs/op

NewFromFloat/ShopSpring        3082452   385.2 ns/op    40 B/op    2 allocs/op
NewFromFloat/Decimal          19404170   59.52 ns/op    24 B/op    1 allocs/op

NewFromInt/ShopSpring       1000000000   0.294 ns/op     0 B/op    0 allocs/op
NewFromInt/Decimal            59144376   19.50 ns/op    16 B/op    1 allocs/op

StringFixed/ShopSpring         3318596   359.6 ns/op   392 B/op   16 allocs/op
StringFixed/Decimal           12386185   96.12 ns/op    40 B/op    2 allocs/op

Abs/ShopSpring                 7241440   163.7 ns/op   200 B/op    7 allocs/op
Abs/Decimal                   13462662   87.47 ns/op    48 B/op    2 allocs/op

Neg/ShopSpring                 5966871   207.3 ns/op   264 B/op    9 allocs/op
Neg/Decimal                   12938702   93.11 ns/op    48 B/op    2 allocs/op

Truncate/ShopSpring            2883730   438.9 ns/op   530 B/op   20 allocs/op
Truncate/Decimal               6978840   173.6 ns/op    72 B/op    3 allocs/op

Ceil/ShopSpring                1833598   693.9 ns/op   872 B/op   30 allocs/op
Ceil/Decimal                   5914952   201.3 ns/op    88 B/op    4 allocs/op

Round/ShopSpring               1446426   860.9 ns/op   920 B/op   34 allocs/op
Round/Decimal                  6058872   188.5 ns/op    72 B/op    3 allocs/op

RoundAwayFromZero/ShopSpring   1870034   642.4 ns/op   872 B/op   30 allocs/op
RoundAwayFromZero/Decimal      6073666   195.0 ns/op    88 B/op    4 allocs/op

RoundTowardToZero/ShopSpring   1839658   653.7 ns/op   872 B/op   30 allocs/op
RoundTowardToZero/Decimal      6925075   171.9 ns/op    72 B/op    3 allocs/op

Floor/ShopSpring               1848147   644.8 ns/op   872 B/op   30 allocs/op
Floor/Decimal                  6430852   187.6 ns/op    72 B/op    3 allocs/op

Shift/ShopSpring               2355440   502.0 ns/op   608 B/op   22 allocs/op
Shift/Decimal                  6324964   187.4 ns/op    88 B/op    4 allocs/op

IntPart/ShopSpring             5684866   197.9 ns/op   264 B/op    9 allocs/op
IntPart/Decimal               12530052   96.62 ns/op    24 B/op    1 allocs/op

IsZero/ShopSpring              7049337   163.5 ns/op   200 B/op    7 allocs/op
IsZero/Decimal                15329292    77.6 ns/op    24 B/op    1 allocs/op

IsInteger/ShopSpring           6387838   164.1 ns/op   200 B/op    7 allocs/op
IsInteger/Decimal             14660770   82.99 ns/op    24 B/op    1 allocs/op

IsPositive/ShopSpring          7241178   163.7 ns/op   200 B/op    7 allocs/op
IsPositive/Decimal            14811858    79.3 ns/op    24 B/op    1 allocs/op

IsNegative/ShopSpring          7227627   165.7 ns/op   200 B/op    7 allocs/op
IsNegative/Decimal            24162208    48.6 ns/op    24 B/op    1 allocs/op

Cmp/ShopSpring                 3219708   363.5 ns/op   424 B/op   15 allocs/op
Cmp/Decimal                    7908812   150.7 ns/op    40 B/op    2 allocs/op

Equal/ShopSpring               3187255   359.3 ns/op   424 B/op   15 allocs/op
Equal/Decimal                  8344328   141.0 ns/op    40 B/op    2 allocs/op

Greater/ShopSpring             3378079   355.7 ns/op   424 B/op   15 allocs/op
Greater/Decimal                8121801   147.6 ns/op    40 B/op    2 allocs/op

Less/ShopSpring                3315786   361.3 ns/op   424 B/op   15 allocs/op
Less/Decimal                   8109493   146.1 ns/op    40 B/op    2 allocs/op

GreaterOrEqual/ShopSpring      3325186   357.7 ns/op   424 B/op   15 allocs/op
GreaterOrEqual/Decimal         8165496   147.7 ns/op    40 B/op    2 allocs/op

LessOrEqual/ShopSpring         3283453   361.6 ns/op   424 B/op   15 allocs/op
LessOrEqual/Decimal            8175444   147.0 ns/op    40 B/op    2 allocs/op

Sign/ShopSpring                7318190   162.8 ns/op   200 B/op    7 allocs/op
Sign/Decimal                  15136862    79.6 ns/op    24 B/op    1 allocs/op

Pow/ShopSpring                  170517    6856 ns/op  6653 B/op  261 allocs/op
Pow/Decimal                     224826    5283 ns/op  2376 B/op  106 allocs/op

Add/ShopSpring                 1877947   646.3 ns/op   728 B/op   28 allocs/op
Add/Decimal                    3236866   366.2 ns/op   136 B/op    6 allocs/op

Sub/ShopSpring                 1838361   649.7 ns/op   744 B/op   29 allocs/op
Sub/Decimal                    3426580   356.9 ns/op   120 B/op    6 allocs/op

Mul/ShopSpring                 4171518   281.8 ns/op   320 B/op   12 allocs/op
Mul/Decimal                    3270048   364.3 ns/op   104 B/op    4 allocs/op

Div/ShopSpring                 2890861   423.2 ns/op   464 B/op   17 allocs/op
Div/Decimal                    1838936   644.4 ns/op   312 B/op   14 allocs/op
```

# Contribution

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## Development

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## API Differences

This section outlines the differences between our decimal implementation and the shopspring/decimal library. We've made several design decisions to improve performance, simplify the API, and provide more intuitive method names.

### In Development

```makefile
Method:
  Mod()
```

### Alternatives:

```makefile
Method:
  Copy()            no need, Decimal is already a copyable structure (string)
  MarshalJSON()     no need, Decimal is already a marshallable structure (string)
  UnmarshalJSON()   no need, Decimal is already a Unmarshallable structure (string)

  Round()           parameter type `int32 -> int`
  Shift()           parameter type `int32 -> int`
  StringFixed()     parameter type `int32 -> int`
  Truncate()        parameter type `int32 -> int`

  GreaterThan()         ->  Greater
  GreaterThanOrEqual()  ->  GreaterOrEqual
  LessThan()            ->  Less
  LessThanOrEqual()     ->  LessOrEqual
  Equals()              ->  Equal
  RoundUp()             ->  RoundAwayFromZero
  RoundDown()           ->  RoundTowardToZero
  RoundFloor()          ->  Floor
  RoundCeil()           ->  Ceil
  Ceil()                ->  Ceil(0)
  Floor()               ->  Floor(0)
```

### Unimplemented:

```makefile
Struct:
  NullDecimal

Variables:
  MarshalJSONWithoutQuotes
  ExpMaxIterations

Function:
  Sum
  Avg
  Max
  Min
  NewFromFloatWithExponent
  NewFromFormattedString
  NewNullDecimal
  RescalePair

Method:
  Atan
  Coefficient
  CoefficientInt64
  Cos
  DivRound
  ExpHullAbrham
  ExpTaylor
  Exponent
  GobDecode
  GobEncode
  InexactFloat64
  NumDigits
  QuoRem
  RoundCash
  Sin
  StringFixedBank
  StringFixedCash
  StringScaled
  Tan
  MarshalBinary
  MarshalText
  UnmarshalBinary
  UnmarshalText
```
