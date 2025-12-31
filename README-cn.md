# Decimal

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)

提供三种固定大小 `Decimal128` / `Decimal256` / `Decimal512` 实现的零分配固定小数十进制库。

## 系统要求

- Go 1.25+（以 `go.mod` 为准）

## 导入

```go
import "github.com/yanun0323/decimal"
```

## 概览

各类型有自己的固定小数尺度与精度：

- `Decimal128`：尺度 `10^19`，整数 **19** 位，小数 **19** 位
- `Decimal256`：尺度 `10^38`，整数 **38** 位，小数 **38** 位
- `Decimal512`：尺度 `10^77`，整数 **77** 位，小数 **77** 位

共通特性：

- **数值模型**：`raw / 10^scale`
- **溢出**：按位宽截断（wrap-around）
- **零值**：可直接使用，表示 `0`
- **不使用 big.Int**：全程固定大小运算

内存布局：

- `Decimal128`：128-bit 二补数（2 x uint64）
- `Decimal256`：256-bit 二补数（4 x uint64）
- `Decimal512`：512-bit 二补数（8 x uint64）

### 精度规则

所有构造与解析都遵守：

- **整数部分**：只保留最低 *N* 位（更高位会被丢弃）
- **小数部分**：只保留最高 *N* 位（更低位会被丢弃）

其中 *N* 为该类型的小数位精度（19/38/77）。

字符串/JSON 解析会**先应用指数位移**，再应用精度规则（截断），最后缩放到固定尺度。

## 预设值

- `Zero128`, `One128`, `Ten128`, `Hundred128`
- `Zero256`, `One256`, `Ten256`, `Hundred256`
- `Zero512`, `One512`, `Ten512`, `Hundred512`

## 构造函数（各类型通用）

将 `XXX` 替换为 `128`、`256` 或 `512`：

- `NewXXX(intPart, decimalPart int64) DecimalXXX`
  - `decimalPart` 视为小数位数字。
  - 例：`New256(123, 45)` = `123.45`。
  - 小数位超过尺度会向 0 截断。
  - 应用精度规则（整数低 *N*、小数高 *N*）。
- `NewXXXFromString(string) (DecimalXXX, error)`
  - 支持前后空白、`_` 分隔、`.` 小数点与 `e/E` 指数。
  - 先应用指数位移，再应用精度规则（截断），最后缩放到固定尺度。
- `NewXXXFromInt(int64) DecimalXXX`
  - 应用精度规则（整数低 *N*）。
- `NewXXXFromFloat(float64) (DecimalXXX, error)`
  - 向 0 截断，`NaN/Inf` 返回错误。
  - 转换后应用精度规则。
- `NewXXXFromBinary([]byte) (DecimalXXX, error)`
  - 固定长度，小端序（见 Binary 章节）。
  - 解码后应用精度规则。
- `NewXXXFromJSON([]byte) (DecimalXXX, error)`
  - 接受 JSON **字符串**或**数字**。
  - 先应用指数位移，再应用精度规则（截断），最后缩放到固定尺度。

## 转换与格式化

- `Int64() (intPart, decimalPart int64)`
  - `decimalPart` 以类型尺度返回（`10^19` / `10^38` / `10^77`）。
- `Float64() float64`
- `String() string`
  - 自动移除小数尾端多余 0。
- `StringFixed(n int) string`
  - `n > scaleDigits` 会截断到类型尺度（19/38/77）。
  - `n <= 0` 只返回整数部分。

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

数位操作规则（scaleDigits = 19/38/77）：

- `n > scaleDigits`：不变
- `n <= -scaleDigits`：返回 0
- `n < 0`：作用在整数位（例如 `Truncate(-1)` 影响十位）

## 超越函数

- `Pow`（仅支持整数次方；指数向 0 截断）
- `Sqrt`, `Exp`, `Log`, `Log2`, `Log10`
  - 无效输入（例如负数开根号、非正数取对数）会返回原值。

## Binary / JSON 编码

- **Binary**：固定长度、小端序
  - `Decimal128`：16 bytes
  - `Decimal256`：32 bytes
  - `Decimal512`：64 bytes
- **JSON**：输出字符串；解析接受字符串或数字

## 数据库集成

- SQL（`database/sql`）
  - 实现 `sql.Scanner` 与 `driver.Valuer`
  - `Scan` 接受 `string`、`[]byte`、`int64`、`float64` 与 `NULL`（NULL 会变成零值）
  - `Value` 返回十进制字符串
- MongoDB Go Driver v2（`go.mongodb.org/mongo-driver/v2`）
  - 实现 `bson.ValueMarshaler` / `bson.ValueUnmarshaler`
  - `Decimal128` 会编码为 BSON Decimal128
  - `Decimal256` / `Decimal512` 会编码为 BSON 字符串

## 泛型接口

包内同时提供编译期限制用的泛型接口：

- `type Decimal[T decimal] interface { ... }`

## 性能

主打零分配与稳定成本。基准测试：

```sh
go test -bench . ./...
```

## 贡献

欢迎提交 PR；重大变更请先开 issue 讨论设计方向。
