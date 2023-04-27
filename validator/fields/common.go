package fields

import (
	"fmt"
	"math/big"

	"github.com/shopspring/decimal"
)

// (substreams - subgraph)/substream * 100.0
func validateErrorPercent(expected, actual *big.Float, errorPercent float64) (bool, float64) {
	numerator := new(big.Float).Add(expected, new(big.Float).Mul(actual, new(big.Float).SetInt64(-1)))
	quotient := new(big.Float).Quo(new(big.Float).Abs(numerator), expected)
	percent := new(big.Float).Mul(quotient, new(big.Float).SetUint64(100))
	v, _ := percent.Float64()

	if percent.Cmp(new(big.Float).SetFloat64(errorPercent)) > 0 {
		return false, v
	}
	return true, v
}

func validTolerance(expected, actual *big.Float, tolerance float64) (bool, float64) {
	tol := new(big.Float).SetFloat64(tolerance)
	dt := new(big.Float).Add(expected, new(big.Float).Mul(actual, new(big.Float).SetInt64(-1)))
	v, _ := dt.Float64()
	if (dt.Cmp(tol) > 0) || dt.Cmp(new(big.Float).Mul(tol, new(big.Float).SetInt64(-1))) < 0 {
		return false, v
	}
	return true, v
}

func validFloatWithShortRound(expected, actual string) bool {
	l1 := len(expected)
	l2 := len(actual)
	if l1 == 0 && l2 == 0 {
		return true
	}
	if l1 == 0 {
		return false
	}
	if l2 == 0 {
		return false
	}

	v1, err := decimal.NewFromString(expected)
	if err != nil {
		panic(fmt.Sprintf("decoding decimal from string %q: %v", expected, err))
	}
	v2, _ := decimal.NewFromString(actual)
	if err != nil {
		panic(fmt.Sprintf("decoding decimal from string %q: %v", actual, err))
	}

	prec := v1.Exponent()
	if v2.Exponent() > prec {
		prec = v2.Exponent()
	}
	prec = -prec

	return v1.Round(prec).Cmp(v2.Round(prec)) == 0
}
