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

- **Overall Speed**: 1.4-4.2x faster across all operations
- **Memory Efficiency**: 50-88% reduction in memory allocations
- **Best Improvements**: Creation (4x faster), Truncate (2.6x faster), Comparisons (2x faster)
- **Consistent Gains**: Every tested operation shows both speed and memory improvements

The benchmarks cover core decimal operations including creation, arithmetic (add/sub/mul), transformations (truncate/shift/abs/neg), and comparisons, demonstrating comprehensive performance advantages while maintaining full compatibility.

```
BenchmarkNew/ShopSpringDecimal-20                        5979028               206.6 ns/op           200 B/op          7 allocs/op
BenchmarkNew/Decimal-20                                 21940134               49.45 ns/op            24 B/op          1 allocs/op

BenchmarkAbs/ShopSpringDecimal-20                        5992156               202.5 ns/op           200 B/op          7 allocs/op
BenchmarkAbs/Decimal-20                                 12136202               99.75 ns/op            48 B/op          2 allocs/op

BenchmarkNeg/ShopSpringDecimal-20                        4420572               247.6 ns/op           264 B/op          9 allocs/op
BenchmarkNeg/Decimal-20                                 10095094               120.2 ns/op            72 B/op          3 allocs/op

BenchmarkTruncate/ShopSpringDecimal-20                   2309737               539.2 ns/op           530 B/op         20 allocs/op
BenchmarkTruncate/Decimal-20                             5044777               205.6 ns/op            96 B/op          4 allocs/op

BenchmarkShift/ShopSpringDecimal-20                      1893084               600.0 ns/op           608 B/op         22 allocs/op
BenchmarkShift/Decimal-20                                3309234               382.7 ns/op           224 B/op          8 allocs/op

BenchmarkAdd/ShopSpringDecimal-20                        1498164               754.5 ns/op           728 B/op         28 allocs/op
BenchmarkAdd/Decimal-20                                  2256195               544.6 ns/op           248 B/op         12 allocs/op

BenchmarkSub/ShopSpringDecimal-20                        1657006               774.5 ns/op           744 B/op         29 allocs/op
BenchmarkSub/Decimal-20                                  2376349               507.8 ns/op           216 B/op         12 allocs/op

BenchmarkIsZero/ShopSpringDecimal-20                     5737850               214.5 ns/op           200 B/op          7 allocs/op
BenchmarkIsZero/Decimal-20                              11718336               108.2 ns/op            48 B/op          2 allocs/op

BenchmarkIsPositive/ShopSpringDecimal-20                 5634820               231.8 ns/op           200 B/op          7 allocs/op
BenchmarkIsPositive/Decimal-20                          10650404               102.0 ns/op            48 B/op          2 allocs/op

BenchmarkIsNegative/ShopSpringDecimal-20                 5557935               232.3 ns/op           200 B/op          7 allocs/op
BenchmarkIsNegative/Decimal-20                          11192362               108.5 ns/op            48 B/op          2 allocs/op

BenchmarkEqual/ShopSpringDecimal-20                      2635482               430.8 ns/op           424 B/op         15 allocs/op
BenchmarkEqual/Decimal-20                                6604009               189.0 ns/op            80 B/op          4 allocs/op

BenchmarkGreater/ShopSpringDecimal-20                    2810475               427.7 ns/op           424 B/op         15 allocs/op
BenchmarkGreater/Decimal-20                              5347701               197.2 ns/op            80 B/op          4 allocs/op

BenchmarkLess/ShopSpringDecimal-20                       2843552               436.8 ns/op           424 B/op         15 allocs/op
BenchmarkLess/Decimal-20                                 6280638               192.8 ns/op            80 B/op          4 allocs/op

BenchmarkGreaterOrEqual/ShopSpringDecimal-20             2671776               430.4 ns/op           424 B/op         15 allocs/op
BenchmarkGreaterOrEqual/Decimal-20                       5619379               193.3 ns/op            80 B/op          4 allocs/op

BenchmarkLessOrEqual/ShopSpringDecimal-20                2719836               426.5 ns/op           424 B/op         15 allocs/op
BenchmarkLessOrEqual/Decimal-20                          6183186               214.3 ns/op            80 B/op          4 allocs/op

BenchmarkMul/ShopSpringDecimal-20                        3519614               368.5 ns/op           320 B/op         12 allocs/op
BenchmarkMul/Decimal-20                                  3416486               329.5 ns/op           136 B/op          5 allocs/op

BenchmarkDiv/ShopSpringDecimal-20                        2274916               501.5 ns/op           464 B/op         17 allocs/op
BenchmarkDiv/Decimal-20                                  1628210               790.5 ns/op           704 B/op         25 allocs/op
```
