package fields

import (
	"fmt"
	"math/big"
)

type Decimal struct {
	v *big.Float

	error     *float64
	tolerance *float64
}

func newDecimal(v *big.Float, opt map[string]interface{}) *Decimal {
	f := &Decimal{v: v}

	f.error, f.tolerance = extractErrorAndTolerance(opt)

	return f
}

func newDecimalFromStr(v string) (Comparable, error) {
	value, ok := new(big.Float).SetString(v)
	if !ok {
		return nil, fmt.Errorf("failed to convert %q to bigfloat", v)
	}
	return &Decimal{v: value}, nil
}

func (f *Decimal) Eql(v Comparable) bool {
	expected := f.v
	actual := v.(*Decimal).v

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

func (f *Decimal) String() string {
	return f.v.String()
}
