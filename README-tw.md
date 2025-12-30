# Decimal256

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)

基於 256-bit 兩補數整數的固定小數（32 位）十進位庫，目標為零記憶體分配與高效能。

## 系統需求

- Go 1.25+（以 `go.mod` 為準）

## 匯入

```go
import "github.com/yanun0323/decimal"
```

## 概覽

Decimal256 為固定 32 位小數的十進位型別：

- **數值模型**：`raw / 10^32`
- **儲存**：256-bit 兩補數（4 x uint64）
- **溢位**：256-bit 截斷（wrap-around）
- **零值**：可直接使用，表示 `0`
- **不使用 big.Int**：全程固定大小運算

## 建構子

- `NewDecimal256(intPart, decimalPart int64) Decimal256`
  - `decimalPart` 視為「小數位數字」。
  - 例：`NewDecimal256(123, 45)` = `123.45`。
  - 超過 32 位小數會被截斷（toward zero）。
- `NewDecimal256FromString(string) (Decimal256, error)`
  - 支援前後空白、`_` 分隔、`.` 小數點與 `e/E` 指數。
  - 小數位超過 32 位會被截斷。
- `NewDecimal256FromInt(int64) Decimal256`
- `NewDecimal256FromFloat(float64) (Decimal256, error)`
  - 向零截斷，`NaN/Inf` 回傳錯誤。
- `NewDecimal256FromBinary([]byte) (Decimal256, error)`
  - 32 bytes，小端序。
- `NewDecimal256FromJSON([]byte) (Decimal256, error)`
  - 接受 JSON **字串**或**數字**。

## 轉換與格式化

- `Int64() (intPart, decimalPart int64)`
  - `decimalPart` 仍為 10^32 的小數尺度。
- `Float64() float64`
- `String() string`
  - 自動移除小數尾端多餘 0。
- `StringFixed(n int) string`
  - `n > 32` 會被截斷為 32；`n <= 0` 只回傳整數部分。

### 零記憶體分配追加

- `AppendBinary(dst []byte) []byte`
- `AppendJSON(dst []byte) []byte`
- `AppendString(dst []byte) []byte`
- `AppendStringFixed(dst []byte, n int) []byte`

以上方法將結果追加到呼叫端提供的 buffer，由呼叫端控制配置行為。

## 檢查

- `IsZero()`, `IsPositive()`, `IsNegative()`, `Sign()`

## 算術與比較

- `Add`, `Sub`, `Mul`, `Div`, `Mod`
  - `Div/Mod` 除以 0 會回傳原值。
- `Equal`, `GreaterThan`, `LessThan`, `GreaterOrEqual`, `LessOrEqual`

## 取整與修正

- `Neg`, `Inv`, `Abs`, `Truncate`, `Shift`
- `Round`（銀行家捨入）
- `RoundAwayFromZero`, `RoundTowardToZero`
- `Ceil`, `Floor`

數位操作規則：

- `n > 32`：不變
- `n <= -32`：回傳 0
- `n < 0`：作用在整數位（例如 `Truncate(-1)` 影響十位）

## 超越函數

- `Pow`（只支援整數次方；指數向 0 截斷）
- `Sqrt`, `Exp`, `Log`, `Log2`, `Log10`
  - 無效輸入（例如負數開根號、非正數取對數）會回傳原值。

## Binary / JSON 編碼

- **Binary**：32 bytes，小端序（4 x uint64）
- **JSON**：輸出字串；解析接受字串或數字

## 效能

主打零記憶體分配與穩定成本。基準測試：

```sh
go test -bench . ./...
```

## 貢獻

歡迎提交 PR；重大變更請先開 issue 討論設計方向。
