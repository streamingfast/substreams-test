package fields

import (
	"fmt"
	"github.com/streamingfast/substreams-test/validator/config"
	"strings"
)

type Array struct {
	values []Comparable
}

type FArrayFactory struct {
	factories []ComparableFactory
}

func (f *FArrayFactory) newArrayFromStr(v string) (Comparable, error) {
	if v == "[]" {
		return &Array{}, nil
	}
	v = strings.TrimPrefix(strings.TrimSuffix(v, "]"), "[")
	chunks := strings.Split(v, ",")
	if len(chunks) != len(f.factories) {
		return nil, fmt.Errorf("unable to parse array %s length does not match expected array", v)
	}
	out := &Array{}
	for i := 0; i < len(f.factories); i++ {
		cleanValue := strings.TrimSpace(strings.TrimLeft(strings.TrimRight(chunks[i], "\""), "\""))
		value, err := f.factories[i](cleanValue)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", cleanValue, err)
		}
		out.values = append(out.values, value)
	}
	return out, nil
}

func (f *Array) Eql(v Comparable) bool {
	in := v.(*Array)

	if len(in.values) != len(f.values) {
		return false
	}

	for i := 0; i < len(f.values); i++ { // [0, 1, 2]
		equals := false

		for j := 0; j < len(in.values); j++ { // [1, 0, 2]
			switch val := f.values[i].(type) {
			case *String:
				_, ok := in.values[j].(*String)

				if !ok { // can't compare string with something else which isn't a string
					return false
				}

				a := newString(strings.ReplaceAll(val.String(), "\"", ""), config.Options{})
				b := newString(strings.ReplaceAll(in.values[j].String(), "\"", ""), config.Options{})

				if strings.Trim(a.v, " ") == strings.Trim(b.v, " ") {
					equals = true
				}

			default:
				if val.Eql(in.values[i]) {
					equals = true
				}
			}
		}

		if !equals {
			return false
		}
	}

	return true
}

func (f *Array) String() string {
	strs := []string{}
	for _, field := range f.values {
		strs = append(strs, field.String())
	}
	return strings.Join(strs, ", ")
}
