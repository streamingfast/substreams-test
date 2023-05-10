package fields

import (
	"fmt"
	"math/big"

	"github.com/streamingfast/substreams-test/validator/config"
)

type Bigint struct {
	config.Options
	v *big.Int
}

func newBigint(v *big.Int, opt config.Options) *Bigint {
	return &Bigint{v: v, Options: opt}
}

func newBigintFromStr(v string, opt config.Options) (Comparable, error) {
	value, ok := new(big.Int).SetString(v, 10)
	if !ok {
		return nil, fmt.Errorf("failed to convert %q to bigint", v)
	}
	return &Bigint{Options: opt, v: value}, nil
}

func (f *Bigint) Eql(v Comparable) bool {
	expected := new(big.Float).SetInt(f.v)
	actual := new(big.Float).SetInt(v.(*Bigint).v)

	if expected.Cmp(actual) == 0 {
		return true
	}

	if f.Tolerance != nil {
		ok, _ := validTolerance(expected, actual, *f.Tolerance)
		return ok
	}

	if f.Error != nil {
		ok, _ := validateErrorPercent(expected, actual, *f.Error)
		return ok
	}

	return false
}

func (f *Bigint) String() string {
	return f.v.String()
}
