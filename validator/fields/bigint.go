package fields

import (
	"fmt"
	"math/big"
)

type Bigint struct {
	v         *big.Int
	error     *float64
	tolerance *float64
}

func newBigint(v *big.Int, opt map[string]interface{}) *Bigint {
	f := &Bigint{v: v}

	f.error, f.tolerance = extractErrorAndTolerance(opt)

	return f
}

func newBigintFromStr(v string) (Comparable, error) {
	value, ok := new(big.Int).SetString(v, 10)
	if !ok {
		return nil, fmt.Errorf("failed to convert %q to bigint", v)
	}
	return &Bigint{v: value}, nil
}

func (f *Bigint) Eql(v Comparable) bool {
	expected := new(big.Float).SetInt(f.v)
	actual := new(big.Float).SetInt(v.(*Bigint).v)

	if f.tolerance != nil {
		ok, _ := validTolerance(expected, actual, *f.tolerance)
		return ok
	}

	if f.error != nil {
		ok, _ := validateErrorPercent(expected, actual, *f.error)
		return ok
	}

	return expected.Cmp(actual) == 0
}

func (f *Bigint) String() string {
	return f.v.String()
}
