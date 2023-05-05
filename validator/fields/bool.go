package fields

import (
	"github.com/streamingfast/substreams-test/validator/config"
)

type Bool struct {
	v bool
}

func newBool(v bool, opt config.Options) *Bool {
	return &Bool{v: v}
}

func newFBoolFromStr(v string, _ config.Options) (Comparable, error) {
	return &Bool{
		v: v == "true",
	}, nil
}

func (f *Bool) Eql(v Comparable) bool {
	return v.(*Bool).v == f.v
}

func (f *Bool) String() string {
	if f.v {
		return "true"
	}
	return "false"
}
