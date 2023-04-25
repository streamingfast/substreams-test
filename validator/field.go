package validator

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/streamingfast/eth-go"

	pbentities "github.com/streamingfast/substreams-test/pb/entity/v1"
)

type Field struct {
	substreamsEntity string
	substreamsField  string
	graphEntity      string
	graphqlField     string
	graphqlJSONPath  string

	objFactory ComparableFactory
	obj        Comparable
}

type ComparableFactory func(string) (Comparable, error)
type Comparable interface {
	eql(v Comparable) bool
	string() string
}

func (v *Validator) newField(substreamsEntity string, field *pbentities.Field) *Field {
	graphqlFieldName := v.getGraphQLFieldName(substreamsEntity, field.Name)
	out := &Field{
		substreamsField: field.Name,
		graphqlField:    graphqlFieldName,
		graphqlJSONPath: fmt.Sprintf("data.%s.%s", normalizeEntityName(substreamsEntity), graphqlFieldName),
	}

	if v.isGraphQLAssociatedField(substreamsEntity, field.Name) {
		out.graphqlField = fmt.Sprintf("%s { id }", graphqlFieldName)
		if v.isGraphQLArrayField(substreamsEntity, field.Name) {
			out.graphqlJSONPath = fmt.Sprintf("data.%s.%s.#.id", normalizeEntityName(substreamsEntity), graphqlFieldName)
		} else {
			out.graphqlJSONPath = fmt.Sprintf("data.%s.%s.id", normalizeEntityName(substreamsEntity), graphqlFieldName)
		}
	}

	fieldOpt := v.getFieldOpt(substreamsEntity, field.Name)

	out.obj, out.objFactory = parseValue(field.NewValue, fieldOpt)

	return out
}

func parseValue(value *pbentities.Value, fieldOpt map[string]interface{}) (Comparable, ComparableFactory) {
	switch newValue := value.Typed.(type) {
	case *pbentities.Value_Int32:
		return newFInt32(newValue.Int32, fieldOpt), newFInt32FromStr

	case *pbentities.Value_Bigdecimal:
		nvalue, _ := new(big.Float).SetString(newValue.Bigdecimal)
		return newFBigDecimal(nvalue, fieldOpt), newFBigDecimalFromStr

	case *pbentities.Value_Bigint:
		nvalue, _ := new(big.Int).SetString(newValue.Bigint, 10)
		return newFBigint(nvalue, fieldOpt), newFBigintFromStr

	case *pbentities.Value_String_:
		return newFString(newValue.String_, fieldOpt), newFStringFromStr

	case *pbentities.Value_Bool:
		return newFBool(newValue.Bool, fieldOpt), newFBoolFromStr

	case *pbentities.Value_Bytes:
		data, err := b64.StdEncoding.DecodeString(newValue.Bytes)
		if err != nil {
			panic(err)
		}

		return newFBytes(data, fieldOpt), newFBytesFromStr
	case *pbentities.Value_Array:
		arr := &FArray{}
		arrFactory := &FArrayFactory{}
		for _, v := range newValue.Array.Value {
			val, factory := parseValue(v, fieldOpt)
			arr.values = append(arr.values, val)
			arrFactory.factories = append(arrFactory.factories, factory)
		}
		return arr, arrFactory.newFArrayFromStr
	default:
		panic(fmt.Errorf("unknown field v value type %T", newValue))
	}

}

type FInt32 struct {
	v int32

	error     *float64
	tolerance *float64
}

func newFInt32(v int32, opt map[string]interface{}) *FInt32 {
	f := &FInt32{v: v}
	f.error, f.tolerance = extractErrorAndTolerance(opt)
	return f

}

func newFInt32FromStr(v string) (Comparable, error) {
	value, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to convert %q to int32: %w", v, err)
	}
	return &FInt32{v: int32(value)}, nil
}

func (f *FInt32) eql(v Comparable) bool {
	expected := new(big.Float).SetInt64(int64(f.v))
	actual := new(big.Float).SetInt64(int64(v.(*FInt32).v))

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

func (f *FInt32) string() string {
	return fmt.Sprintf("%d", f.v)
}

type FBigdecimal struct {
	v *big.Float

	error     *float64
	tolerance *float64
}

func newFBigDecimal(v *big.Float, opt map[string]interface{}) *FBigdecimal {
	f := &FBigdecimal{v: v}

	f.error, f.tolerance = extractErrorAndTolerance(opt)

	return f
}

func newFBigDecimalFromStr(v string) (Comparable, error) {
	value, ok := new(big.Float).SetString(v)
	if !ok {
		return nil, fmt.Errorf("failed to convert %q to bigfloat", v)
	}
	return &FBigdecimal{v: value}, nil
}

func (f *FBigdecimal) eql(v Comparable) bool {
	expected := f.v
	actual := v.(*FBigdecimal).v

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

func (f *FBigdecimal) string() string {
	return f.v.String()
}

type FBigint struct {
	v         *big.Int
	error     *float64
	tolerance *float64
}

func newFBigint(v *big.Int, opt map[string]interface{}) *FBigint {
	f := &FBigint{v: v}

	f.error, f.tolerance = extractErrorAndTolerance(opt)

	return f
}

func newFBigintFromStr(v string) (Comparable, error) {
	value, ok := new(big.Int).SetString(v, 10)
	if !ok {
		return nil, fmt.Errorf("failed to convert %q to bigint", v)
	}
	return &FBigint{v: value}, nil
}

func (f *FBigint) eql(v Comparable) bool {
	expected := new(big.Float).SetInt(f.v)
	actual := new(big.Float).SetInt(v.(*FBigint).v)

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

func (f *FBigint) string() string {
	return f.v.String()
}

type FString struct {
	v       string
	process func(v string) string
}

func newFString(v string, opt map[string]interface{}) *FString {
	return &FString{
		v: v,

		process: func(v string) string {
			return v
		},
	}
}

func newFStringFromStr(v string) (Comparable, error) {
	return &FString{
		v: v,
	}, nil
}

func (f *FString) eql(v Comparable) bool {
	return f.v == f.process(v.(*FString).v)
}

func (f *FString) string() string {
	return f.v
}

type FBool struct {
	v bool
}

func newFBool(v bool, opt map[string]interface{}) *FBool {
	return &FBool{v: v}
}

func newFBoolFromStr(v string) (Comparable, error) {
	return &FBool{
		v: v == "true",
	}, nil
}

func (f *FBool) eql(v Comparable) bool {
	return v.(*FBool).v == f.v
}

func (f *FBool) string() string {
	if f.v {
		return "true"
	}
	return "false"
}

type FBytes struct {
	v []byte
}

func newFBytes(v []byte, opt map[string]interface{}) *FBytes {
	return &FBytes{v: v}
}

func newFBytesFromStr(v string) (Comparable, error) {
	var data []byte
	var err error
	if strings.HasPrefix(v, "0x") {
		if data, err = hex.DecodeString(eth.SanitizeHex(v)); err != nil {
			return nil, fmt.Errorf("failed to convert %s to byte array: %w", v, err)
		}
	}
	return &FBytes{v: data}, nil
}

func (f *FBytes) eql(v Comparable) bool {
	return bytes.Compare(f.v, v.(*FBytes).v) == 0
}

func (f *FBytes) string() string {
	return hex.EncodeToString(f.v)
}

type FArray struct {
	values []Comparable
}

type FArrayFactory struct {
	factories []ComparableFactory
}

func (f *FArrayFactory) newFArrayFromStr(v string) (Comparable, error) {
	if v == "[]" {
		return &FArray{}, nil
	}
	v = strings.TrimPrefix(strings.TrimSuffix(v, "]"), "[")
	chunks := strings.Split(v, ",")
	if len(chunks) != len(f.factories) {
		return nil, fmt.Errorf("unable to parse array %s length does not match expected array", v)
	}
	out := &FArray{}
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

func (f *FArray) eql(v Comparable) bool {
	in := v.(*FArray)

	if len(in.values) != len(f.values) {
		return false
	}
	for i := 0; i < len(in.values); i++ {
		if !f.values[i].eql(in.values[i]) {
			return false
		}
	}
	return true
}

func (f *FArray) string() string {
	strs := []string{}
	for _, field := range f.values {
		strs = append(strs, field.string())
	}
	return strings.Join(strs, ", ")
}

func extractErrorAndTolerance(opts map[string]interface{}) (*float64, *float64) {

	errpercInt, errpercIntOk := opts["error"]
	toleranceInt, toleranceIntOk := opts["tolerance"]

	if errpercIntOk && toleranceIntOk {
		panic("error and tolerance are mutually exclusive when comparing numerical values")
	}

	if errpercIntOk {
		v := errpercInt.(float64)
		return &v, nil
	}

	if toleranceIntOk {
		v := toleranceInt.(float64)
		return nil, &v
	}
	return nil, nil
}

// (substreams - subgraph)/substream * 100.0
func validateErrorPercent(expected, actual *big.Float, errorPercent float64) (bool, float64) {
	numerator := new(big.Float).Add(expected, new(big.Float).Mul(actual, new(big.Float).SetInt64(-1)))
	quotient := new(big.Float).Quo(new(big.Float).Abs(numerator), expected)
	percent := new(big.Float).Mul(quotient, new(big.Float).SetUint64(100))
	v, _ := percent.Float64()

	if percent.Cmp(new(big.Float).SetFloat64(errorPercent)) > 0 {
		return false, v
	}
	return true, v
}

func validTolerance(expected, actual *big.Float, tolerance float64) (bool, float64) {
	tol := new(big.Float).SetFloat64(tolerance)
	dt := new(big.Float).Add(expected, new(big.Float).Mul(actual, new(big.Float).SetInt64(-1)))
	v, _ := dt.Float64()
	if (dt.Cmp(tol) > 0) || dt.Cmp(new(big.Float).Mul(tol, new(big.Float).SetInt64(-1))) < 0 {
		return false, v
	}
	return true, v

}
