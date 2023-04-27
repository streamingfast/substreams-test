package fields

import (
	"fmt"
	"math/big"

	"github.com/streamingfast/substreams-test/validator/config"
)

type Decimal struct {
	config.Options
	v   *big.Float
	str string
}

func newDecimal(v *big.Float, opt config.Options) *Decimal {
	return &Decimal{v: v, Options: opt}
}

func newDecimalFromStr(v string) (Comparable, error) {
	value, ok := new(big.Float).SetString(v)
	if !ok {
		return nil, fmt.Errorf("failed to convert %q to bigfloat", v)
	}
	return &Decimal{v: value, str: v}, nil
}

func (f *Decimal) Eql(v Comparable) bool {
	expected := f.v
	actual := v.(*Decimal).v

	if f.Tolerance != nil {
		ok, _ := validTolerance(expected, actual, *f.Tolerance)
		return ok
	}

	if f.Round != "" {
		switch f.Round {
		case "shortest":
			expected := f.str
			actual := v.(*Decimal).str
			return validFloatWithShortRound(expected, actual)
		default:
			panic(fmt.Sprintf("unsupported round mode %q", f.Round))
		}
	}

	if f.Error != nil {
		ok, _ := validateErrorPercent(expected, actual, *f.Error)
		return ok
	}

	return expected.Cmp(actual) == 0
}

func (f *Decimal) String() string {
	return f.str
}
