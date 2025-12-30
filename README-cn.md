# Decimal256

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)

基于 256-bit 二补数整数的固定小数（32 位）十进制库，目标为零分配与高性能。

## 系统要求

- Go 1.25+（以 `go.mod` 为准）

## 导入

```go
import "github.com/yanun0323/decimal"
```

## 概览

Decimal256 为固定 32 位小数的十进制类型：

- **数值模型**：`raw / 10^32`
- **存储**：256-bit 二补数（4 x uint64）
- **溢出**：256-bit 截断（wrap-around）
- **零值**：可直接使用，表示 `0`
- **不使用 big.Int**：全程固定大小运算

## 构造函数

- `NewDecimal256(intPart, decimalPart int64) Decimal256`
  - `decimalPart` 视为“小数位数字”。
  - 例：`NewDecimal256(123, 45)` = `123.45`。
  - 超过 32 位小数会被截断（toward zero）。
- `NewDecimal256FromString(string) (Decimal256, error)`
  - 支持前后空白、`_` 分隔、`.` 小数点与 `e/E` 指数。
  - 小数位超过 32 位会被截断。
- `NewDecimal256FromInt(int64) Decimal256`
- `NewDecimal256FromFloat(float64) (Decimal256, error)`
  - 向零截断，`NaN/Inf` 返回错误。
- `NewDecimal256FromBinary([]byte) (Decimal256, error)`
  - 32 bytes，小端序。
- `NewDecimal256FromJSON([]byte) (Decimal256, error)`
  - 接受 JSON **字符串**或**数字**。

## 转换与格式化

- `Int64() (intPart, decimalPart int64)`
  - `decimalPart` 仍为 10^32 的小数尺度。
- `Float64() float64`
- `String() string`
  - 自动移除小数尾端多余 0。
- `StringFixed(n int) string`
  - `n > 32` 会被截断为 32；`n <= 0` 只返回整数部分。

### 零分配追加

- `AppendBinary(dst []byte) []byte`
- `AppendJSON(dst []byte) []byte`
- `AppendString(dst []byte) []byte`
- `AppendStringFixed(dst []byte, n int) []byte`

以上方法将结果追加到调用方提供的 buffer，由调用方控制分配行为。

## 检查

- `IsZero()`, `IsPositive()`, `IsNegative()`, `Sign()`

## 算术与比较

- `Add`, `Sub`, `Mul`, `Div`, `Mod`
  - `Div/Mod` 除以 0 会返回原值。
- `Equal`, `GreaterThan`, `LessThan`, `GreaterOrEqual`, `LessOrEqual`

## 取整与修正

- `Neg`, `Inv`, `Abs`, `Truncate`, `Shift`
- `Round`（银行家舍入）
- `RoundAwayFromZero`, `RoundTowardToZero`
- `Ceil`, `Floor`

数位操作规则：

- `n > 32`：不变
- `n <= -32`：返回 0
- `n < 0`：作用在整数位（例如 `Truncate(-1)` 影响十位）

## 超越函数

- `Pow`（仅支持整数次方；指数向 0 截断）
- `Sqrt`, `Exp`, `Log`, `Log2`, `Log10`
  - 无效输入（例如负数开根号、非正数取对数）会返回原值。

## Binary / JSON 编码

- **Binary**：32 bytes，小端序（4 x uint64）
- **JSON**：输出字符串；解析接受字符串或数字

## 性能

主打零分配与稳定成本。基准测试：

```sh
go test -bench . ./...
```

## 贡献

欢迎提交 PR；重大变更请先开 issue 讨论设计方向。
