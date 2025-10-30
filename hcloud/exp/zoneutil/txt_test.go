package zoneutil

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestCaseTXT struct {
	desc    string
	raw     string
	encoded string
}

var (
	manyA        = strings.Repeat("a", 255)
	fewB         = strings.Repeat("b", 10)
	testCasesTXT = []TestCaseTXT{
		{
			desc:    "empty",
			raw:     ``,
			encoded: ``,
		},
		{
			desc:    "empty quotes",
			raw:     `""`,
			encoded: `"\"\""`,
		},
		{
			desc:    "small",
			raw:     `hello world`,
			encoded: `"hello world"`,
		},
		{
			desc:    "long",
			raw:     manyA + fewB,
			encoded: fmt.Sprintf(`"%s" "%s"`, manyA, fewB),
		},
		{
			desc:    "new line",
			raw:     "hello\nworld",
			encoded: `"hello` + "\n" + `world"`,
		},
		{
			desc:    "quotes even",
			raw:     `hello "world"`,
			encoded: `"hello \"world\""`,
		},
		{
			desc:    "quotes odd",
			raw:     `hello "world`,
			encoded: `"hello \"world"`,
		},
	}
)

func TestFormatTXTValue(t *testing.T) {
	for _, testCase := range testCasesTXT {
		t.Run(testCase.desc, func(t *testing.T) {
			got := FormatTXTRecord(testCase.raw)
			require.Equal(t, testCase.encoded, got)
		})
	}
}

func TestParseTXTValue(t *testing.T) {
	for _, testCase := range testCasesTXT {
		t.Run(testCase.desc, func(t *testing.T) {
			got := ParseTXTRecord(testCase.encoded)
			require.Equal(t, testCase.raw, got)
		})
	}
}
