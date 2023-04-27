package fields

import (
	"math/big"
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

func validFloatWithPrecision(expected, actual *big.Float, precision int) bool {
	roundedExpected := new(big.Float).SetPrec(uint(precision)).SetMode(big.ToZero).Copy(expected)
	roundedActual := new(big.Float).SetPrec(uint(precision)).SetMode(big.ToZero).Copy(actual)
	return roundedExpected.Cmp(roundedActual) == 0
}

func validFloatWithShortRound(expected, actual *big.Float) bool {
	precision := expected.Prec()
	if actPrec := actual.Prec(); actPrec < precision {
		precision = actPrec
	}
	return validFloatWithPrecision(expected, actual, int(precision))
}
