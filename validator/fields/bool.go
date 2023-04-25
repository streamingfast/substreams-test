package fields

type Bool struct {
	v bool
}

func newBool(v bool, opt map[string]interface{}) *Bool {
	return &Bool{v: v}
}

func newFBoolFromStr(v string) (Comparable, error) {
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
