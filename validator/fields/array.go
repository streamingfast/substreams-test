package fields

import (
	"fmt"
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
		cleanValue := strings.TrimLeft(strings.TrimRight(chunks[i], "\""), "\"")
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
	for i := 0; i < len(in.values); i++ {
		if !f.values[i].Eql(in.values[i]) {
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
