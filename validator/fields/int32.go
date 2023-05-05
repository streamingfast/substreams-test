package fields

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/streamingfast/substreams-test/validator/config"
)

type Int32 struct {
	config.Options
	v int32
}

func newInt32(v int32, opt config.Options) *Int32 {
	return &Int32{
		Options: opt,
		v:       v,
	}
}

func newInt32FromStr(v string, opt config.Options) (Comparable, error) {
	value, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to convert %q to int32: %w", v, err)
	}
	return &Int32{Options: opt, v: int32(value)}, nil
}

func (f *Int32) Eql(v Comparable) bool {
	expected := new(big.Float).SetInt64(int64(f.v))
	actual := new(big.Float).SetInt64(int64(v.(*Int32).v))

	if f.Tolerance != nil {
		ok, _ := validTolerance(expected, actual, *f.Tolerance)
		return ok
	}

	if f.Error != nil {
		ok, _ := validateErrorPercent(expected, actual, *f.Error)
		return ok
	}

	return expected.Cmp(actual) == 0

}

func (f *Int32) String() string {
	return fmt.Sprintf("%d", f.v)
}
