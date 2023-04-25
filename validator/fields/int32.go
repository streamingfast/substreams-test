package fields

import (
	"fmt"
	"math/big"
	"strconv"
)

type Int32 struct {
	v int32

	error     *float64
	tolerance *float64
}

func newInt32(v int32, opt map[string]interface{}) *Int32 {
	f := &Int32{v: v}
	f.error, f.tolerance = extractErrorAndTolerance(opt)
	return f

}

func newInt32FromStr(v string) (Comparable, error) {
	value, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to convert %q to int32: %w", v, err)
	}
	return &Int32{v: int32(value)}, nil
}

func (f *Int32) Eql(v Comparable) bool {
	expected := new(big.Float).SetInt64(int64(f.v))
	actual := new(big.Float).SetInt64(int64(v.(*Int32).v))

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

func (f *Int32) String() string {
	return fmt.Sprintf("%d", f.v)
}
