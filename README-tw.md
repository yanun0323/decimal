# Decimal

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)

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
- 完全相容 [shopspring/decimal](https://github.com/shopspring/decimal) API - 所有函數都實作為支援相同介面
- 任何差異或未實作的功能都記錄在「[API 差異](README.md#api-differences)」章節中

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

詳細的效能基準測試請參考 [English README](README.md#benchmark)。

# 貢獻

歡迎貢獻！請隨時提交 Pull Request。對於重大變更，請先開啟 issue 討論您想要變更的內容。

## 開發

1. Fork 此儲存庫
2. 建立您的功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的變更 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 開啟 Pull Request

## API 差異

詳細的 API 差異說明請參考 [English README](README.md#api-differences)。
