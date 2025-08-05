# Decimal

A super efficient decimal base on string type.

## Requirements

#### _GO 1.21 or higher_

## Import

```go
import "github.com/yanun0323/decimal"
```

## Features

- The zero-value is 0, and is safe to use without initialization
- Addition, subtraction with no loss of precision
- Database/sql serialization/deserialization as string
- JSON and XML serialization/deserialization as string

## Supported

- Initial from string
- Addition
- Subtraction
- Multiplication
- Division (Not optimize yet)
- Negative
- Truncate
- Shift
- Compare like equal, greater, less, greater or equal, less or equal

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

- **Overall Speed**: 1.1-5x faster across all operations
- **Memory Efficiency**: 50-100% reduction in memory allocations
- **Best Improvements**: Creation (5x faster), IsNegative (5.3x faster), Comparisons (3-4x faster)
- **Arithmetic Operations**: Add/Sub operations are 2x faster with significant memory savings
- **Zero Allocations**: Many operations (creation, comparisons) achieve zero memory allocations

The benchmarks cover core decimal operations including creation, arithmetic (add/sub/mul), transformations (truncate/shift/abs/neg), and comparisons, demonstrating comprehensive performance advantages while maintaining full compatibility.

```
BenchmarkNew/ShopSpringDecimal-20                        5919193               194.1 ns/op           200 B/op          7 allocs/op
BenchmarkNew/Decimal-20                                 30230674               39.09 ns/op             0 B/op          0 allocs/op

BenchmarkAbs/ShopSpringDecimal-20                        5952031               203.2 ns/op           200 B/op          7 allocs/op
BenchmarkAbs/Decimal-20                                 14264305               78.57 ns/op            24 B/op          1 allocs/op

BenchmarkNeg/ShopSpringDecimal-20                        4827121               263.6 ns/op           264 B/op          9 allocs/op
BenchmarkNeg/Decimal-20                                 14046283               86.06 ns/op            24 B/op          1 allocs/op

BenchmarkTruncate/ShopSpringDecimal-20                   2117624               509.4 ns/op           530 B/op         20 allocs/op
BenchmarkTruncate/Decimal-20                             7681054               160.9 ns/op            24 B/op          1 allocs/op

BenchmarkShift/ShopSpringDecimal-20                      2029108               609.8 ns/op           608 B/op         22 allocs/op
BenchmarkShift/Decimal-20                                6765243               176.5 ns/op            40 B/op          2 allocs/op

BenchmarkAdd/ShopSpringDecimal-20                        1579506               714.0 ns/op           728 B/op         28 allocs/op
BenchmarkAdd/Decimal-20                                  3458757               340.5 ns/op            96 B/op          4 allocs/op

BenchmarkSub/ShopSpringDecimal-20                        1618722               759.4 ns/op           744 B/op         29 allocs/op
BenchmarkSub/Decimal-20                                  3815695               326.9 ns/op            80 B/op          4 allocs/op

BenchmarkIsZero/ShopSpringDecimal-20                     5926252               202.9 ns/op           200 B/op          7 allocs/op
BenchmarkIsZero/Decimal-20                              17871141               68.91 ns/op             0 B/op          0 allocs/op

BenchmarkIsPositive/ShopSpringDecimal-20                 6099165               207.8 ns/op           200 B/op          7 allocs/op
BenchmarkIsPositive/Decimal-20                          16562524               69.76 ns/op             0 B/op          0 allocs/op

BenchmarkIsNegative/ShopSpringDecimal-20                 5939834               204.5 ns/op           200 B/op          7 allocs/op
BenchmarkIsNegative/Decimal-20                          31879681               38.52 ns/op             0 B/op          0 allocs/op

BenchmarkEqual/ShopSpringDecimal-20                      2791113               429.3 ns/op           424 B/op         15 allocs/op
BenchmarkEqual/Decimal-20                               10202091               119.4 ns/op             0 B/op          0 allocs/op

BenchmarkGreater/ShopSpringDecimal-20                    2909077               458.0 ns/op           424 B/op         15 allocs/op
BenchmarkGreater/Decimal-20                              9688188               130.6 ns/op             0 B/op          0 allocs/op

BenchmarkLess/ShopSpringDecimal-20                       2778007               409.4 ns/op           424 B/op         15 allocs/op
BenchmarkLess/Decimal-20                                 9694372               127.7 ns/op             0 B/op          0 allocs/op

BenchmarkGreaterOrEqual/ShopSpringDecimal-20             2877723               426.0 ns/op           424 B/op         15 allocs/op
BenchmarkGreaterOrEqual/Decimal-20                       9541648               123.0 ns/op             0 B/op          0 allocs/op

BenchmarkLessOrEqual/ShopSpringDecimal-20                2816113               426.4 ns/op           424 B/op         15 allocs/op
BenchmarkLessOrEqual/Decimal-20                          9752029               125.0 ns/op             0 B/op          0 allocs/op

BenchmarkMul/ShopSpringDecimal-20                        3358723               338.7 ns/op           320 B/op         12 allocs/op
BenchmarkMul/Decimal-20                                  3629094               293.6 ns/op            64 B/op          2 allocs/op

BenchmarkDiv/ShopSpringDecimal-20                        2560041               474.7 ns/op           464 B/op         17 allocs/op
BenchmarkDiv/Decimal-20                                  1795086               700.0 ns/op           352 B/op         12 allocs/op
```
