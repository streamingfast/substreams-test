package fields

import (
	"github.com/streamingfast/substreams-test/validator/config"
)

type String struct {
	v       string
	process func(v string) string
}

func newString(v string, opt config.Options) *String {
	return &String{
		v: v,

		process: func(v string) string {
			return v
		},
	}
}

func newStringFromStr(v string) (Comparable, error) {
	return &String{
		v: v,
	}, nil
}

func (f *String) Eql(v Comparable) bool {
	return f.v == f.process(v.(*String).v)
}

func (f *String) String() string {
	return f.v
}
