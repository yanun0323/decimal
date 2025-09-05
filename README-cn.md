# Decimal

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)

基于字符串类型的超高效率、内存优化十进制数字运算库。

## 系统要求

#### _GO 1.21 或更高版本_

## 导入

```go
import "github.com/yanun0323/decimal"
```

## 特色功能

- 零值默认为 0，无需初始化即可安全使用
- 内存优化，在维持高性能的同时拥有极低的内存占用
- 加减法运算无精度损失
- 支持 Database/sql 序列化/反序列化
- 支持 JSON 和 XML 以字符串形式序列化/反序列化
- 完全兼容 [shopspring/decimal](https://github.com/shopspring/decimal) API - 所有函数都实现为支持相同接口
- 任何差异或未实现的功能都记录在「[API 差异](README.md#api-differences)」章节中

## 支持功能

- 初始化：支持从字符串、int、int32、float、float64、big.Int 创建
- 加法运算
- 减法运算
- 乘法运算
- 除法运算
- 负数运算
- 截断
- 位移
- 比较：等于、大于、小于、大于等于、小于等于
- 四舍五入：一般四舍五入、无条件进位、无条件舍去、银行家四舍五入、远离零值四舍五入、趋向零值四舍五入

## 使用方法

```go
// 创建 decimal
zero := decimal.Zero()

d1, err := decimal.New("100,000.555")

d2 := decimal.Require("50_000.05")

// 运算操作
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

// 比较运算
d1.IsZero()             // false
d1.IsPositive()         // true
d1.IsNegative()         // true

d1.Equal(d2)            // false
d1.Greater(d2)          // true
d1.Less(d2)             // false
d1.GreaterOrEqual(d2)   // true
d1.LessOrEqual(d2)      // false


// 方法链式调用
result := d1.Sub(d2).Shift(-5).Add(d1).Truncate(3).String()
```

## 性能基准测试

详细的性能基准测试请参考 [English README](README.md#benchmark)。

# 贡献

欢迎贡献！请随时提交 Pull Request。对于重大变更，请先开启 issue 讨论您想要变更的内容。

## 开发

1. Fork 此仓库
2. 创建您的功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的变更 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启 Pull Request

## API 差异

详细的 API 差异说明请参考 [English README](README.md#api-differences)。
