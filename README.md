# Decimal
A string base decimal.

## Supported go versions
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
- Truncate
- Shift

## Unsupported Operation
- Multiplication
- Division

## Usage
```go
// create a decimal with error
d1, err := decimal.New("100,000.555")

d2 := decimal.Require("50_000.05")

// operation
add := d1.Add(d2).String()
println(add)            // 150000.605

sub := d1.Sub(d2).String()
println(sub)            // 50000.505

shift := d1.Shift(-2).String()
println(shift)          // 1000.00555

truncate := d1.Truncate(1).String()
println(truncate)       // 100000.5


// method chain
result := d1.Sub(d2).Shift(-5).Add(d1).Truncate(3).String()
```

## Benchmark
Compare to [github.com/shopspring/decimal](https://github.com/shopspring/decimal)
> Memory usage less than half of `shopspring/decimal`

![Alt text](./benchmark.png)