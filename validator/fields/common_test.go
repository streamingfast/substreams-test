package fields

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/assert"
)

func TestValidator_errorPercent(t *testing.T) {
	bf := func(in string) *big.Float {
		e, _, err := big.ParseFloat(in, 10, 64, big.ToNearestAway)
		if err != nil {
			require.NoError(t, err)
		}
		return e
	}
	tests := []struct {
		expect             *big.Float
		actual             *big.Float
		errorPercent       float64
		expectBool         bool
		expecterrorPercent float64
	}{
		{
			expect:             new(big.Float).SetFloat64(2.24),
			actual:             new(big.Float).SetFloat64(2.12),
			errorPercent:       7.3,
			expectBool:         true,
			expecterrorPercent: 5.3571428571,
		},
		{
			expect:             new(big.Float).SetFloat64(2.24),
			actual:             new(big.Float).SetFloat64(2.12),
			errorPercent:       5.3,
			expectBool:         false,
			expecterrorPercent: 5.3571428571,
		},
		{
			expect:             bf("0.000311147009695436285092442318800308911778113924614013850621487019379910004202707027884960798574336598"),
			actual:             bf("0.0003111470096954362733861181240203676"),
			errorPercent:       0.0000000000001,
			expectBool:         true,
			expecterrorPercent: 1e-14,
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			match, errPerc := validateErrorPercent(test.expect, test.actual, test.errorPercent)
			assert.Equal(t, test.expectBool, match)
			assert.InDelta(t, test.expecterrorPercent, errPerc, 0.0001)
		})
	}

}

func TestValidator_validTolerance(t *testing.T) {

	tests := []struct {
		expect          *big.Float
		actual          *big.Float
		tolerance       float64
		expectBool      bool
		expectTolerance float64
	}{
		{
			expect:          new(big.Float).SetFloat64(1.232932),
			actual:          new(big.Float).SetFloat64(1.232723),
			tolerance:       0.00001,
			expectBool:      false,
			expectTolerance: 0.00020899999999990370725,
		},
		{
			expect:          new(big.Float).SetFloat64(1.232932),
			actual:          new(big.Float).SetFloat64(1.232723),
			tolerance:       0.001,
			expectBool:      true,
			expectTolerance: 0.00020899999999990370725,
		},
		{
			expect:          new(big.Float).SetFloat64(1.232932),
			actual:          new(big.Float).SetFloat64(1.232),
			tolerance:       0.001,
			expectBool:      true,
			expectTolerance: 0.0009319999999999329,
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			match, tolerance := validTolerance(test.expect, test.actual, test.tolerance)
			assert.Equal(t, test.expectBool, match)
			assert.InDelta(t, test.expectTolerance, tolerance, 0.000001)
		})
	}

}

func TestSomething(t *testing.T) {
	assert.True(t, validFloatWithShortRound("0.999", "0.9992"))
	assert.False(t, validFloatWithShortRound("0.999", "0.9999"))
}
