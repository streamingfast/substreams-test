package fields

import (
	"math/big"
	"strings"
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

var _10b = big.NewInt(10)

func validFloatWithPrecision(expected, actual *big.Float, precision uint64) bool {
	bigDecimals := new(big.Int).Exp(_10b, big.NewInt(int64(precision)), nil)
	expectedWhole, _ := new(big.Float).Mul(expected, new(big.Float).SetInt(bigDecimals)).Int(nil)
	actualWhole, _ := new(big.Float).Mul(actual, new(big.Float).SetInt(bigDecimals)).Int(nil)
	return expectedWhole.Cmp(actualWhole) == 0
}

func validFloatWithShortRound(expected, actual *big.Float) bool {
	precision := decimalCount(expected)
	if actPrec := decimalCount(actual); actPrec < precision {
		precision = actPrec
	}
	return validFloatWithPrecision(expected, actual, precision)
}

func decimalCount(v *big.Float) uint64 {
	chunks := strings.Split(v.String(), ".")
	if len(chunks) == 1 {
		return 0
	}
	return uint64(len(chunks[1]))
}
