# Decimal

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)
[![日本語](https://img.shields.io/badge/日本語-クリック-青)](README-ja.md)

基於字串型別的超高效率、記憶體優化十進制數字運算庫。

## 系統需求

#### _GO 1.21 或更高版本_

## 匯入

```go
import "github.com/yanun0323/decimal"
```

## 特色功能

- 零值預設為 0，無需初始化即可安全使用
- 記憶體優化，在維持高效能的同時擁有極低的記憶體佔用
- 加減法運算無精度損失
- 支援 Database/sql 序列化/反序列化
- 支援 JSON 和 XML 以字串形式序列化/反序列化

## 支援功能

- 初始化：支援從字串、int、int32、float、float64、big.Int 建立
- 加法運算
- 減法運算
- 乘法運算
- 除法運算
- 負數運算
- 截斷
- 位移
- 比較：等於、大於、小於、大於等於、小於等於
- 四捨五入：一般四捨五入、無條件進位、無條件捨去、銀行家四捨五入、遠離零值四捨五入、趨向零值四捨五入

## 使用方法

```go
// 建立 decimal
zero := decimal.Zero()

d1, err := decimal.New("100,000.555")

d2 := decimal.Require("50_000.05")

// 運算操作
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

// 比較運算
d1.IsZero()             // false
d1.IsPositive()         // true
d1.IsNegative()         // true

d1.Equal(d2)            // false
d1.Greater(d2)          // true
d1.Less(d2)             // false
d1.GreaterOrEqual(d2)   // true
d1.LessOrEqual(d2)      // false


// 方法鏈式呼叫
result := d1.Sub(d2).Shift(-5).Add(d1).Truncate(3).String()
```

## 效能基準測試

與 [github.com/shopspring/decimal](https://github.com/shopspring/decimal) 比較

- **整體速度**：所有操作都快 1.9-6.5 倍
- **記憶體效率**：記憶體配置減少 70-88%
- **建立操作**：New 快 3 倍，NewFromFloat 快 6.5 倍，且大幅節省記憶體
- **四捨五入操作**：所有四捨五入方法都顯示 3-4 倍的效能改進，並大幅減少記憶體使用

這些基準測試在保持與 shopspring/decimal 完全 API 相容性的同時，展現了在建立、算術、轉換和比較方面一致的效能優勢。

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

IsZero/ShopSpring              7049337   163.5 ns/op   200 B/op    7 allocs/op
IsZero/Decimal                15329292    77.6 ns/op    24 B/op    1 allocs/op

IsPositive/ShopSpring          7241178   163.7 ns/op   200 B/op    7 allocs/op
IsPositive/Decimal            14811858    79.3 ns/op    24 B/op    1 allocs/op

IsNegative/ShopSpring          7227627   165.7 ns/op   200 B/op    7 allocs/op
IsNegative/Decimal            24162208    48.6 ns/op    24 B/op    1 allocs/op

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

Add/ShopSpring                 1877947   646.3 ns/op   728 B/op   28 allocs/op
Add/Decimal                    3236866   366.2 ns/op   136 B/op    6 allocs/op

Sub/ShopSpring                 1838361   649.7 ns/op   744 B/op   29 allocs/op
Sub/Decimal                    3426580   356.9 ns/op   120 B/op    6 allocs/op

Mul/ShopSpring                 4171518   281.8 ns/op   320 B/op   12 allocs/op
Mul/Decimal                    3270048   364.3 ns/op   104 B/op    4 allocs/op

Div/ShopSpring                 2890861   423.2 ns/op   464 B/op   17 allocs/op
Div/Decimal                    1838936   644.4 ns/op   312 B/op   14 allocs/op
```
