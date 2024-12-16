# Decimal

A high memory optimized decimal library for Go.

Check the [benchmark](./README.md#benchmark) for more details.

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
- Division _unoptimized_
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

_All the benchmark start from the `NewFromString` function, please check the [benchmark.go](./decimal_bench_test.go) for more details._

### Memory

> smaller is better

| Operation                  | _yanun0323/decimal_ | _shopspring/decimal_ |
| -------------------------- | ------------------- | -------------------- |
| Addition                   | 264 B/op            | 728 B/op             |
| Subtract                   | 248 B/op            | 744 B/op             |
| Multiply                   | 136 B/op            | 320 B/op             |
| Division</br>_Unoptimized_ | 704 B/op            | 484 B/op             |
| Shift                      | 248 B/op            | 608 B/op             |
| Abs                        | 24 B/op             | 200 B/op             |
| Neg                        | 24 B/op             | 264 B/op             |
| Truncate                   | 48 B/op             | 530 B/op             |
| IsZero                     | 24 B/op             | 200 B/op             |
| IsPositive                 | 24 B/op             | 200 B/op             |
| IsNegative                 | 24 B/op             | 200 B/op             |
| Equal                      | 40 B/op             | 424 B/op             |
| Greater                    | 40 B/op             | 424 B/op             |
| Less                       | 40 B/op             | 424 B/op             |

### Speed

> smaller is better

| Operation                  | _yanun0323/decimal_ | _shopspring/decimal_ |
| -------------------------- | ------------------- | -------------------- |
| Addition                   | 625.6 ns/op         | 685.3 ns/op          |
| Subtract                   | 637.1 ns/op         | 681.3 ns/op          |
| Multiply                   | 767.0 ns/op         | 314.5 ns/op          |
| Division</br>_Unoptimized_ | 1069 ns/op          | 467.3 ns/op          |
| Shift                      | 661.4 ns/op         | 568.3 ns/op          |
| Abs                        | 209.2 ns/op         | 191.8 ns/op          |
| Neg                        | 229.1 ns/op         | 270.1 ns/op          |
| Truncate                   | 433.4 ns/op         | 484.9 ns/op          |
| IsZero                     | 205.4 ns/op         | 185.6 ns/op          |
| IsPositive                 | 209.8 ns/op         | 195.9 ns/op          |
| IsNegative                 | 204.1 ns/op         | 194.9 ns/op          |
| Equal                      | 352.5 ns/op         | 406.4 ns/op          |
| Greater                    | 367.8 ns/op         | 401.5 ns/op          |
| Less                       | 366.5 ns/op         | 392.9 ns/op          |

### Conclusion

> CPU usage almost same or less than `shopspring/decimal`
> (except multiplication/division)
>
> Memory usage less than half of `shopspring/decimal`
> (except multiplication/division)

### Reference

![Benchmark](./benchmark.png)
