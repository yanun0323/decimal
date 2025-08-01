package decimal

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DecimalSuite struct {
	suite.Suite
}

func TestDecimalSuite(t *testing.T) {
	suite.Run(t, new(DecimalSuite))
}

func (su *DecimalSuite) SetupSuite() {}

func (su *DecimalSuite) TestPressure() {
	d := decimal.RequireFromString("123.123").Truncate(5)
	su.Equal("123.123", d.String())

	b := strings.Repeat("2", 100_000_000)
	a := strings.Repeat("1", 100_000_000)

	bd, err := New(b)
	su.Require().NoError(err)
	ad, err := New(a)
	su.Require().NoError(err)

	res := bd.Sub(ad)
	su.Require().NoError(err)
	su.Equal(a, res.String())
}

func (su *DecimalSuite) TestNewDecimal() {
	testCases := []struct {
		desc     string
		input    string
		hasError bool
		expected string
	}{
		{
			desc:     "Good Without Dropping",
			input:    "-100,000.000,000",
			expected: "-100000",
		},
		{
			desc:     "Good With Dropping",
			input:    "+.000_000_000",
			expected: "0",
		},
		{
			desc:     "Bad Symbol",
			input:    "-100,000.0+00,00-0",
			hasError: true,
		},
		{
			desc:     "Duplicate Dot",
			input:    "-100,000.000.000",
			hasError: true,
		},
		{
			desc:     "Empty String",
			input:    "",
			expected: "0",
		},
		{
			desc:     "Empty String",
			input:    "0",
			expected: "0",
		},
		{
			desc:     "Dot String",
			input:    ".",
			expected: "0",
		},
		{
			desc:     "Dot Zero String",
			input:    ".0",
			expected: "0",
		},
		{
			desc:     "Zero Dot Zero String",
			input:    "0.0",
			expected: "0",
		},
		{
			desc:     "Zero Dot String",
			input:    "0.",
			expected: "0",
		},
		{
			desc:     "Empty String After Dropping",
			input:    "+,,,",
			hasError: true,
		},
		{
			desc:     "Invalid Symbol",
			input:    "&10000000",
			hasError: true,
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			d, err := New(tc.input)
			if tc.hasError {
				su.Require().Error(err)
				return
			}

			su.Require().NoError(err, tc.desc)
			su.Equal(tc.expected, d.String(), tc.desc)
		})
	}
}

func (su *DecimalSuite) TestZeroValue() {
	var d Decimal
	d2, err := New("123")
	su.Require().NoError(err)

	su.NotPanics(func() {
		su.Equal("123", d.Add(d2).String())
		su.Equal("-123", d.Sub(d2).String())
		su.Equal("123", d2.Add(d).String())
		su.Equal("123", d2.Sub(d).String())
		su.Equal("0", d.Truncate(5).String())
		su.Equal("0", d.Shift(5).String())
	})

}

func (su *DecimalSuite) TestAbs() {
	testCases := []struct {
		desc     string
		d        string
		expected string
	}{
		{
			d:        "0",
			expected: "0",
		},
		{
			d:        "123.456",
			expected: "123.456",
		},
		{
			d:        "-123.456",
			expected: "123.456",
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			d, err := New(tc.d)
			su.Require().NoError(err)
			su.Equal(tc.expected, d.Abs().String())
		})
	}
}

func (su *DecimalSuite) TestNeg() {
	testCases := []struct {
		desc     string
		d        string
		expected string
	}{
		{
			d:        "0",
			expected: "0",
		},
		{
			d:        "123.456",
			expected: "-123.456",
		},
		{
			d:        "-123.456",
			expected: "123.456",
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			d, err := New(tc.d)
			su.Require().NoError(err)
			su.Equal(tc.expected, d.Neg().String())
		})
	}
}

func (su *DecimalSuite) TestTruncate() {
	testCases := []struct {
		desc     string
		input    string
		truncate int
		expected string
	}{
		{
			desc:     "Normal",
			input:    "123.123",
			truncate: 0,
			expected: "123",
		},
		{
			desc:     "Normal",
			input:    "123",
			truncate: 0,
			expected: "123",
		},
		{
			desc:     "Normal",
			input:    "123.000",
			truncate: 0,
			expected: "123",
		},
		{
			desc:     "Normal",
			input:    "123.123",
			truncate: 1,
			expected: "123.1",
		},
		{
			desc:     "No Need Truncate",
			input:    "123.123",
			truncate: 5,
			expected: "123.123",
		},
		{
			desc:     "Negative Truncate",
			input:    "123.123",
			truncate: -1,
			expected: "120",
		},
		{
			desc:     "Negative Overflow Truncate",
			input:    "123.123",
			truncate: -3,
			expected: "0",
		},
		{
			desc:     "Natural Number",
			input:    "123",
			truncate: 3,
			expected: "123",
		},
		{
			desc:     "Zero Decimal",
			input:    "123.000",
			truncate: 2,
			expected: "123",
		},
		{
			desc:     "Zero",
			input:    "0",
			truncate: 2,
			expected: "0",
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			d, err := New(tc.input)
			su.Require().NoError(err, tc.desc)
			su.Equal(tc.expected, d.Truncate(tc.truncate).String(), tc.desc)
		})
	}
}

func (su *DecimalSuite) TestFindOrInsertDecimalPoint() {
	testCases := []struct {
		desc          string
		input         string
		expectedNum   string
		expectedIndex int
	}{
		{
			desc:          "A",
			input:         "123.456",
			expectedNum:   "123.456",
			expectedIndex: 3,
		},
		{
			desc:          "B",
			input:         "123123",
			expectedNum:   "123123.",
			expectedIndex: 6,
		},
		{
			desc:          "C",
			input:         ".123",
			expectedNum:   ".123",
			expectedIndex: 0,
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			result, index := findOrInsertDecimalPoint([]byte(tc.input))
			su.Equal(tc.expectedNum, string(result), tc.desc)
			su.Equal(tc.expectedIndex, index, tc.desc)
		})
	}
}

func (su *DecimalSuite) TestUnsignedAdd() {
	testCases := []struct {
		desc           string
		base           string
		addition       string
		expectedResult string
	}{
		{
			desc:           "A",
			base:           "123.456",
			addition:       "123123.456456",
			expectedResult: "123246.912456",
		},
		{
			desc:           "B",
			base:           "123123.0",
			addition:       "123.456456",
			expectedResult: "123246.456456",
		},
		{
			desc:           "C",
			base:           "123",
			addition:       "0.8888",
			expectedResult: "123.8888",
		},
		{
			desc:           "D",
			base:           "123123.544",
			addition:       "900123.456",
			expectedResult: "1023247",
		},
		{
			desc:           "E",
			base:           "123123",
			addition:       "123",
			expectedResult: "123246",
		},
		{
			desc:           "F",
			base:           "0.00001",
			addition:       "0.02",
			expectedResult: "0.02001",
		},
		{
			desc:           "G",
			base:           "123000",
			addition:       "877000",
			expectedResult: "1000000",
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			result := unsignedAdd([]byte(tc.base), []byte(tc.addition))
			su.Equal(tc.expectedResult, string(result), tc.desc)
		})
	}
}

func (su *DecimalSuite) TestUnsignedSub() {
	testCases := []struct {
		desc           string
		base           string
		subtraction    string
		expectedResult string
	}{
		{
			desc:           "A",
			base:           "123.456",
			subtraction:    "123123.456456",
			expectedResult: "-123000.000456",
		},
		{
			desc:           "B",
			base:           "123123.000000",
			subtraction:    "000123.456456",
			expectedResult: "122999.543544",
		},
		{
			desc:           "C",
			base:           "123",
			subtraction:    "0.8888",
			expectedResult: "122.1112",
		},
		{
			desc:           "D",
			base:           "123123.544",
			subtraction:    "123.456",
			expectedResult: "123000.088",
		},
		{
			desc:           "E",
			base:           "123123",
			subtraction:    "123",
			expectedResult: "123000",
		},
		{
			desc:           "F",
			base:           "0.00001",
			subtraction:    "0.02",
			expectedResult: "-0.01999",
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			result := unsignedSub([]byte(tc.base), []byte(tc.subtraction))
			su.Equal(tc.expectedResult, string(result), tc.desc)
		})
	}
}

func (su *DecimalSuite) TestCleanZero() {
	testCases := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "Suffix Zero",
			input:    "000.00001000",
			expected: "0.00001",
		},
		{
			desc:     "All Zero",
			input:    "00000000",
			expected: "0",
		},
		{
			desc:     "All Zero With Decimal Point",
			input:    "0000.00000",
			expected: "0",
		},
		{
			desc:     "Decimal Point In The End",
			input:    "000888000.",
			expected: "888000",
		},
		{
			desc:     "Decimal Point In The Beginning",
			input:    ".000888000",
			expected: "0.000888",
		},
		{
			desc:     "No Decimal Point",
			input:    "123",
			expected: "123",
		},
		{
			desc:     "No Decimal Point",
			input:    "1",
			expected: "1",
		},
		{
			desc:     "1.0",
			input:    "1.0",
			expected: "1",
		},
		{
			desc:     ".1",
			input:    ".1",
			expected: "0.1",
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			result := tidy([]byte(tc.input))
			su.Equal(tc.expected, string(result), tc.desc)
		})
	}
}

func (su *DecimalSuite) TestDecimalAdd() {
	testCases := []struct {
		desc     string
		base     Decimal
		addition Decimal
		expected Decimal
	}{
		{
			desc:     "Natural Number",
			base:     "123456789",
			addition: "987654321",
			expected: "1111111110",
		},
		{
			desc:     "Two Decimal Number",
			base:     "12345.6789",
			addition: "98765.4321",
			expected: "111111.111",
		},
		{
			desc:     "Base Decimal Number",
			base:     "12345",
			addition: "98765.4321",
			expected: "111110.4321",
		},
		{
			desc:     "Delta Decimal Number",
			base:     "12345.6789",
			addition: "98765",
			expected: "111110.6789",
		},
		{
			desc:     "Natural Base Negative Number",
			base:     "-123456789",
			addition: "987654321",
			expected: "864197532",
		},
		{
			desc:     "Natural Delta Negative Number",
			base:     "123456789",
			addition: "-987654321",
			expected: "-864197532",
		},
		{
			desc:     "Natural Twd Negative Number",
			base:     "-123456789",
			addition: "-987654321",
			expected: "-1111111110",
		},
		{
			desc:     "Both Positive",
			base:     "222.222",
			addition: "111.111",
			expected: "333.333",
		},
		{
			desc:     "Base Negative",
			base:     "-222.222",
			addition: "111.111",
			expected: "-111.111",
		},
		{
			desc:     "Addition Negative",
			base:     "222.222",
			addition: "-111.111",
			expected: "111.111",
		},
		{
			desc:     "Both Negative",
			base:     "-222.222",
			addition: "-111.111",
			expected: "-333.333",
		},
		{
			desc:     "Both Zero",
			base:     "0",
			addition: "0",
			expected: "0",
		},
		{
			desc:     "Positive and Negative Zero",
			base:     "0",
			addition: "-0",
			expected: "0",
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			result := tc.base.Add(tc.addition)
			su.Equal(tc.expected, result, tc.desc)
		})
	}
}

func (su *DecimalSuite) TestDecimalSub() {
	testCases := []struct {
		desc     string
		base     Decimal
		addition Decimal
		expected Decimal
	}{
		{
			desc:     "Natural Number",
			base:     "123456789",
			addition: "987654321",
			expected: "-864197532",
		},
		{
			desc:     "Two Decimal Number",
			base:     "12345.6789",
			addition: "98765.4321",
			expected: "-86419.7532",
		},
		{
			desc:     "Base Decimal Number",
			base:     "12345",
			addition: "98765.4321",
			expected: "-86420.4321",
		},
		{
			desc:     "Delta Decimal Number",
			base:     "12345.6789",
			addition: "98765",
			expected: "-86419.3211",
		},
		{
			desc:     "Natural Base Negative Number",
			base:     "-123456789",
			addition: "987654321",
			expected: "-1111111110",
		},
		{
			desc:     "Natural Delta Negative Number",
			base:     "123456789",
			addition: "-987654321",
			expected: "1111111110",
		},
		{
			desc:     "Natural Twd Negative Number",
			base:     "-123456789",
			addition: "-987654321",
			expected: "864197532",
		},
		{
			desc:     "Both Positive",
			base:     "222.222",
			addition: "111.111",
			expected: "111.111",
		},
		{
			desc:     "Base Negative",
			base:     "-222.222",
			addition: "111.111",
			expected: "-333.333",
		},
		{
			desc:     "Addition Negative",
			base:     "222.222",
			addition: "-111.111",
			expected: "333.333",
		},
		{
			desc:     "Both Negative",
			base:     "-222.222",
			addition: "-111.111",
			expected: "-111.111",
		},
		{
			desc:     "Both Zero",
			base:     "0",
			addition: "0",
			expected: "0",
		},
		{
			desc:     "Positive and Negative Zero",
			base:     "0",
			addition: "-0",
			expected: "0",
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			result := tc.base.Sub(tc.addition)
			su.Equal(tc.expected, result, tc.desc)
		})
	}
}

func TestShift(t *testing.T) {
	testCases := []struct {
		input    string
		shift    int
		expected string
	}{
		{
			input:    "100",
			shift:    8,
			expected: "10000000000",
		},
		{
			input:    "100.123456789",
			shift:    8,
			expected: "10012345678.9",
		},
		{
			input:    "100.12345678",
			shift:    8,
			expected: "10012345678",
		},
		{
			input:    "0.123456789",
			shift:    8,
			expected: "12345678.9",
		},
		{
			input:    "0.12345",
			shift:    8,
			expected: "12345000",
		},
		{
			input:    "-100",
			shift:    8,
			expected: "-10000000000",
		},
		{
			input:    "-100.123456789",
			shift:    8,
			expected: "-10012345678.9",
		},
		{
			input:    "-100.12345678",
			shift:    8,
			expected: "-10012345678",
		},
		{
			input:    "-0.123456789",
			shift:    8,
			expected: "-12345678.9",
		},
		{
			input:    "-0.12345",
			shift:    8,
			expected: "-12345000",
		},
		{
			input:    "10000000000",
			shift:    -8,
			expected: "100",
		},
		{
			input:    "1",
			shift:    -8,
			expected: "0.00000001",
		},
		{
			input:    "10012345678.9",
			shift:    -8,
			expected: "100.123456789",
		},
		{
			input:    "10012345678",
			shift:    -8,
			expected: "100.12345678",
		},
		{
			input:    "12345678.9",
			shift:    -8,
			expected: "0.123456789",
		},
		{
			input:    "123456789",
			shift:    -8,
			expected: "1.23456789",
		},
		{
			input:    "12345000",
			shift:    -8,
			expected: "0.12345",
		},
		{
			input:    "-10000000000",
			shift:    -8,
			expected: "-100",
		},
		{
			input:    "-1",
			shift:    -8,
			expected: "-0.00000001",
		},
		{
			input:    "-1",
			shift:    -2,
			expected: "-0.01",
		},
		{
			input:    "-10012345678.9",
			shift:    -6,
			expected: "-10012.3456789",
		},
		{
			input:    "-10012345678",
			shift:    -8,
			expected: "-100.12345678",
		},
		{
			input:    "-12345678.9",
			shift:    -8,
			expected: "-0.123456789",
		},
		{
			input:    "-123456789",
			shift:    -8,
			expected: "-1.23456789",
		},
		{
			input:    "-12345000",
			shift:    -8,
			expected: "-0.12345",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			t.Log(tc.input)
			d, err := New(tc.input)
			require.NoError(t, err)
			shifted := d.Shift(tc.shift)
			assert.Equal(t, tc.expected, shifted.String(), "input: %s, shift: %d, expected: %s, got: %s", tc.input, tc.shift, tc.expected, shifted.String())
		})
	}
}

func (su *DecimalSuite) TestJsonSupport() {
	d := Require("123456.123456")

	b, err := json.Marshal(d)
	su.Require().NoError(err)
	su.Equal("\"123456.123456\"", string(b))

	var dd Decimal
	su.Require().NoError(json.Unmarshal(b, &dd))
	su.Equal("123456.123456", dd.String())
}

func (su *DecimalSuite) TestIsZero() {
	testCases := []struct {
		desc     string
		d        Decimal
		expected bool
	}{
		{
			d:        "0",
			expected: true,
		},
		{
			d:        "0.0",
			expected: true,
		},
		{
			d:        "00.0000",
			expected: true,
		},
		{
			d:        ".0",
			expected: true,
		},
		{
			d:        ".0",
			expected: true,
		},
		{
			d:        "",
			expected: true,
		},
		{
			d:        "123",
			expected: false,
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			su.Equal(tc.expected, tc.d.IsZero())
		})
	}
}

func (su *DecimalSuite) TestIsPositive() {
	testCases := []struct {
		desc     string
		d        Decimal
		expected bool
	}{
		{
			d:        "123456.88",
			expected: true,
		},
		{
			d:        "0",
			expected: false,
		},
		{
			d:        "0",
			expected: false,
		},
		{
			d:        "-1",
			expected: false,
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			su.Equal(tc.expected, tc.d.IsPositive())
		})
	}
}

func (su *DecimalSuite) TestIsNegative() {
	testCases := []struct {
		desc     string
		d        Decimal
		expected bool
	}{
		{
			d:        "-123456.88",
			expected: true,
		},
		{
			d:        "0",
			expected: false,
		},
		{
			d:        "1",
			expected: false,
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			su.Equal(tc.expected, tc.d.IsNegative())
		})
	}
}

func (su *DecimalSuite) TestEqual() {
	testCases := []struct {
		desc   string
		d1, d2 Decimal
	}{
		{
			d1: "123456.88",
			d2: "123456.88",
		},
		{
			d1: "0",
			d2: "0",
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			su.True(tc.d1.Equal(tc.d2))
		})
	}
}

func (su *DecimalSuite) TestGreater() {
	testCases := []struct {
		desc     string
		d1, d2   Decimal
		expected bool
	}{
		{
			d1:       "123456.89",
			d2:       "123456.88",
			expected: true,
		},
		{
			d1:       "123456.89",
			d2:       "12345.888899",
			expected: true,
		},
		{
			d1:       "0.00001",
			d2:       "0",
			expected: true,
		},
		{
			d1:       "0.00001",
			d2:       "0.00001",
			expected: false,
		},
		{
			d1:       "0",
			d2:       "0",
			expected: false,
		},
		{
			d1:       "12345.888899",
			d2:       "123456.89",
			expected: false,
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			su.Equal(tc.expected, tc.d1.Greater(tc.d2), "%s > %s ?", tc.d1, tc.d2)
		})
	}
}

func (su *DecimalSuite) TestLess() {
	testCases := []struct {
		desc     string
		d1, d2   Decimal
		expected bool
	}{
		{
			d1:       "123456.88",
			d2:       "123456.89",
			expected: true,
		},
		{
			d1:       "12345.888899",
			d2:       "123456.89",
			expected: true,
		},
		{
			d1:       "0",
			d2:       "0.00001",
			expected: true,
		},
		{
			d1:       "0.00001",
			d2:       "0.00001",
			expected: false,
		},
		{
			d1:       "0",
			d2:       "0",
			expected: false,
		},
		{
			d1:       "123456.89",
			d2:       "12345.888899",
			expected: false,
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc, tc.d1, tc.d2)
			su.Equal(tc.expected, tc.d1.Less(tc.d2), "%s < %s ?", tc.d1, tc.d2)
		})
	}
}

func (su *DecimalSuite) TestGreaterOrEqual() {
	testCases := []struct {
		desc     string
		d1, d2   Decimal
		expected bool
	}{
		{
			d1:       "123456.89",
			d2:       "123456.88",
			expected: true,
		},
		{
			d1:       "123456.89",
			d2:       "12345.888899",
			expected: true,
		},
		{
			d1:       "0.00001",
			d2:       "0",
			expected: true,
		},
		{
			d1:       "0.00001",
			d2:       "0.00001",
			expected: true,
		},
		{
			d1:       "0",
			d2:       "0",
			expected: true,
		},
		{
			d1:       "12345.888899",
			d2:       "123456.89",
			expected: false,
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			su.Equal(tc.expected, tc.d1.GreaterOrEqual(tc.d2), "%s >= %s ?", tc.d1, tc.d2)
		})
	}
}

func (su *DecimalSuite) TestLessOrEqual() {
	testCases := []struct {
		desc     string
		d1, d2   Decimal
		expected bool
	}{
		{
			d1:       "123456.88",
			d2:       "123456.89",
			expected: true,
		},
		{
			d1:       "12345.888899",
			d2:       "123456.89",
			expected: true,
		},
		{
			d1:       "0",
			d2:       "0.00001",
			expected: true,
		},
		{
			d1:       "0.00001",
			d2:       "0.00001",
			expected: true,
		},
		{
			d1:       "0",
			d2:       "0",
			expected: true,
		},
		{
			d1:       "123456.89",
			d2:       "12345.888899",
			expected: false,
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			su.Equal(tc.expected, tc.d1.LessOrEqual(tc.d2), "%s <= %s ?", tc.d1, tc.d2)
		})
	}
}

func (su *DecimalSuite) TestMultiplyPureNumber() {
	d1 := []byte("12345")
	d2 := []byte("5648935")

	result := string(multiplyPureNumber(d1, d2, false))
	su.Equal("069736102575", result)
}

func (su *DecimalSuite) TestRemoveDecimalPoint() {
	result, right := removeDecimalPoint([]byte("123.45678"))
	su.Equal("12345678", string(result))
	su.Equal(5, right)
}

func (su *DecimalSuite) TestMul() {
	testCases := []struct {
		desc     string
		d1, d2   Decimal
		expected string
	}{
		{
			d1:       "12345",
			d2:       "5648935",
			expected: "69736102575",
		},
		{
			d1:       "-12345",
			d2:       "5648935",
			expected: "-69736102575",
		},
		{
			d1:       "123.45",
			d2:       "-5648.935",
			expected: "-697361.02575",
		},
		{
			d1:       "-12345",
			d2:       "-5648.935",
			expected: "69736102.575",
		},
		{
			d1:       "123.45",
			d2:       "56.48935",
			expected: "6973.6102575",
		},
		{
			d1:       "100.45",
			d2:       "1000",
			expected: "100450",
		},
		{
			d1:       "100",
			d2:       "1000",
			expected: "100000",
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			result := tc.d1.Mul(tc.d2)
			su.Equal(tc.expected, result.String(), "%s * %s = %s", tc.d1, tc.d2, tc.expected)
		})
	}
}

func (su *DecimalSuite) TestDiv() {
	testCases := []struct {
		desc     string
		d1, d2   Decimal
		expected string
	}{
		{
			d1:       "123",
			d2:       "123",
			expected: "1",
		},
		{
			d1:       "123123123",
			d2:       "123",
			expected: "1001001",
		},
		{
			d1:       "123123.123",
			d2:       "0.123",
			expected: "1001001",
		},
		{
			d1:       "0.000123",
			d2:       "0.123",
			expected: "0.001",
		},
		{
			d1:       "-0.000123123123",
			d2:       "0.123",
			expected: "-0.001001001",
		},
		{
			d1:       "-0.000123",
			d2:       "0.123",
			expected: "-0.001",
		},
		{
			d1:       "-0.000123123123",
			d2:       "-0.123",
			expected: "0.001001001",
		},
		{
			d1:       "-0.000123",
			d2:       "-0.123",
			expected: "0.001",
		},
		{
			d1:       "0.000123123123",
			d2:       "0.123",
			expected: "0.001001001",
		},
		{
			d1:       "-123123123",
			d2:       "123",
			expected: "-1001001",
		},
		{
			d1:       "123123123",
			d2:       "-123",
			expected: "-1001001",
		},
		{
			d1:       "10000",
			d2:       "300",
			expected: "33.3333333333333333",
		},
		{
			d1:       "10000",
			d2:       "1",
			expected: "10000",
		},
		{
			d1:       "10000",
			d2:       "10",
			expected: "1000",
		},
		{
			d1:       "10000",
			d2:       "1000000",
			expected: "0.01",
		},
		{
			d1:       "10000",
			d2:       "1000",
			expected: "10",
		},
		{
			d1:       "84658.4",
			d2:       "333.452",
			expected: "253.8848170051461679",
		},
	}

	for _, tc := range testCases {
		su.T().Run(tc.desc, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("%s, %s, %s, %s", tc.desc, tc.d1, tc.d2, r)
				}
			}()
			result := tc.d1.Div(tc.d2)
			su.Equal(tc.expected, result.String(),
				"input: %s, %s, expected: %s, result: %s",
				tc.d1, tc.d2, tc.expected, result)
		})
	}
}
