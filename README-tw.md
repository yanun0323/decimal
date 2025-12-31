# Decimal

[![English](https://img.shields.io/badge/English-Click-yellow)](README.md)
[![繁體中文](https://img.shields.io/badge/繁體中文-點擊查看-orange)](README-tw.md)
[![简体中文](https://img.shields.io/badge/简体中文-点击查看-orange)](README-cn.md)

提供三種固定大小 `Decimal128` / `Decimal256` / `Decimal512` 實作的零記憶體分配固定小數十進位庫。


## 系統需求

- Go 1.25+（以 `go.mod` 為準）

## 匯入

```go
import "github.com/yanun0323/decimal"
```

## 概覽

各型別有自己的固定小數尺度與精度：

- `Decimal128`：尺度 `10^19`，整數 **19** 位，小數 **19** 位
- `Decimal256`：尺度 `10^38`，整數 **38** 位，小數 **38** 位
- `Decimal512`：尺度 `10^77`，整數 **77** 位，小數 **77** 位

共通特性：

- **數值模型**：`raw / 10^scale`
- **溢位**：依位元寬度截斷（wrap-around）
- **零值**：可直接使用，表示 `0`
- **不使用 big.Int**：全程固定大小運算

記憶體配置：

- `Decimal128`：128-bit 兩補數（2 x uint64）
- `Decimal256`：256-bit 兩補數（4 x uint64）
- `Decimal512`：512-bit 兩補數（8 x uint64）

### 精度規則

所有建構子與解析都遵守：

- **整數部分**：只保留最低 *N* 位（較高位數會被丟棄）
- **小數部分**：只保留最高 *N* 位（較低位數會被丟棄）

其中 *N* 為該型別的小數位精度（19/38/77）。

字串/JSON 解析會**先套用指數位移**，再套用精度規則（截斷），最後才縮放到固定尺度。

## 預設值

- `Zero128`, `One128`, `Ten128`, `Hundred128`
- `Zero256`, `One256`, `Ten256`, `Hundred256`
- `Zero512`, `One512`, `Ten512`, `Hundred512`

## 建構子（各型別共用）

將 `XXX` 取代為 `128`、`256` 或 `512`：

- `NewXXX(intPart, decimalPart int64) DecimalXXX`
  - `decimalPart` 視為小數位數字。
  - 例：`New256(123, 45)` = `123.45`。
  - 小數位超過尺度會被向 0 截斷。
  - 套用精度規則（整數低 *N*、小數高 *N*）。
- `NewXXXFromString(string) (DecimalXXX, error)`
  - 支援前後空白、`_` 分隔、`.` 小數點與 `e/E` 指數。
  - 先套用指數位移，再套用精度規則（截斷），最後縮放到固定尺度。
- `NewXXXFromInt(int64) DecimalXXX`
  - 套用精度規則（整數低 *N*）。
- `NewXXXFromFloat(float64) (DecimalXXX, error)`
  - 向 0 截斷，`NaN/Inf` 回傳錯誤。
  - 轉換後套用精度規則。
- `NewXXXFromBinary([]byte) (DecimalXXX, error)`
  - 固定長度，小端序（見 Binary 章節）。
  - 解碼後套用精度規則。
- `NewXXXFromJSON([]byte) (DecimalXXX, error)`
  - 接受 JSON **字串**或**數字**。
  - 先套用指數位移，再套用精度規則（截斷），最後縮放到固定尺度。

## 轉換與格式化

- `Int64() (intPart, decimalPart int64)`
  - `decimalPart` 以型別尺度回傳（`10^19` / `10^38` / `10^77`）。
- `Float64() float64`
- `String() string`
  - 自動移除小數尾端多餘 0。
- `StringFixed(n int) string`
  - `n > scaleDigits` 會截斷到型別尺度（19/38/77）。
  - `n <= 0` 只回傳整數部分。

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

數位操作規則（scaleDigits = 19/38/77）：

- `n > scaleDigits`：不變
- `n <= -scaleDigits`：回傳 0
- `n < 0`：作用在整數位（例如 `Truncate(-1)` 影響十位）

## 超越函數

- `Pow`（只支援整數次方；指數向 0 截斷）
- `Sqrt`, `Exp`, `Log`, `Log2`, `Log10`
  - 無效輸入（例如負數開根號、非正數取對數）會回傳原值。

## Binary / JSON 編碼

- **Binary**：固定長度、小端序
  - `Decimal128`：16 bytes
  - `Decimal256`：32 bytes
  - `Decimal512`：64 bytes
- **JSON**：輸出字串；解析接受字串或數字

## 資料庫整合

- SQL（`database/sql`）
  - 實作 `sql.Scanner` 與 `driver.Valuer`
  - `Scan` 接受 `string`、`[]byte`、`int64`、`float64` 與 `NULL`（NULL 會變成零值）
  - `Value` 回傳十進位字串
- MongoDB Go Driver v2（`go.mongodb.org/mongo-driver/v2`）
  - 實作 `bson.ValueMarshaler` / `bson.ValueUnmarshaler`
  - `Decimal128` 會編碼為 BSON Decimal128
  - `Decimal256` / `Decimal512` 會編碼為 BSON 字串

## 泛型介面

套件同時提供編譯期限制用的泛型介面：

- `type Decimal[T decimal] interface { ... }`

## 效能

主打零記憶體分配與穩定成本。基準測試：

```sh
go test -bench . ./...
```

## 貢獻

歡迎提交 PR；重大變更請先開 issue 討論設計方向。
