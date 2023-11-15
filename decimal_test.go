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

	bd, err := NewDecimal(b)
	su.Require().NoError(err)
	ad, err := NewDecimal(a)
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
			hasError: true,
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
			d, err := NewDecimal(tc.input)
			if tc.hasError {
				su.Require().Error(err)
				return
			}

			su.Require().NoError(err, tc.desc)
			su.Equal(tc.expected, d.String(), tc.desc)
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
			expected: "123.123",
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
			d, err := NewDecimal(tc.input)
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
			result := cleanZero([]byte(tc.input))
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

func Test_Shift(t *testing.T) {
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
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			t.Log(tc.input)
			result, err := NewDecimal(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result.Shift(tc.shift).String())
		})
	}
}

func (su *DecimalSuite) TestJsonSupport() {
	d := RequireDecimal("123456.123456")

	b, err := json.Marshal(d)
	su.Require().NoError(err)
	su.Equal("\"123456.123456\"", string(b))

	var dd Decimal
	su.Require().NoError(json.Unmarshal(b, &dd))
	su.Equal("123456.123456", dd.String())
}
