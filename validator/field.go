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

	obj Comparable
}
type Comparable interface {
	eql(v string) (bool, error)
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
		out.graphqlJSONPath = fmt.Sprintf("data.%s.%s.id", normalizeEntityName(substreamsEntity), graphqlFieldName)
	}

	fieldOpt := v.getFieldOpt(substreamsEntity, field.Name)

	out.obj = parseValue(field.NewValue, fieldOpt)

	return out
}

func parseValue(value *pbentities.Value, fieldOpt map[string]interface{}) Comparable {
	switch newValue := value.Typed.(type) {
	case *pbentities.Value_Int32:
		return newFInt32(newValue.Int32, fieldOpt)

	case *pbentities.Value_Bigdecimal:
		nvalue, _ := new(big.Float).SetString(newValue.Bigdecimal)
		return newFBigDecimal(nvalue, fieldOpt)

	case *pbentities.Value_Bigint:
		nvalue, _ := new(big.Int).SetString(newValue.Bigint, 10)
		return newFBigint(nvalue, fieldOpt)

	case *pbentities.Value_String_:
		return newFString(newValue.String_, fieldOpt)

	case *pbentities.Value_Bool:
		return newFBool(newValue.Bool, fieldOpt)

	case *pbentities.Value_Bytes:
		data, err := b64.StdEncoding.DecodeString(newValue.Bytes)
		if err != nil {
			panic(err)
		}

		return newFBytes(data, fieldOpt)
	case *pbentities.Value_Array:
		var out []Comparable
		for _, v := range newValue.Array.Value {
			out = append(out, parseValue(v, fieldOpt))
		}
		return newFArray(out)
	default:
		panic(fmt.Errorf("unknown field v value type %T", newValue))
	}

}

type FInt32 struct {
	v int32
}

func newFInt32(v int32, opt map[string]interface{}) *FInt32 {
	return &FInt32{v: v}
}

func (f *FInt32) eql(v string) (bool, error) {
	value, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return false, fmt.Errorf("failed to convert %q to int32: %w", v, err)
	}
	if f.v == int32(value) {
		return true, nil
	}
	return false, nil
}

func (f *FInt32) string() string {
	return fmt.Sprintf("%d", f.v)
}

type FBigdecimal struct {
	v *big.Float
}

func newFBigDecimal(v *big.Float, opt map[string]interface{}) *FBigdecimal {
	return &FBigdecimal{v: v}
}

func (f *FBigdecimal) eql(v string) (bool, error) {
	value, ok := new(big.Float).SetString(v)
	if !ok {
		return false, fmt.Errorf("failed to convert %q to bigfloat", v)
	}
	return f.v.Cmp(value) == 0, nil
}

func (f *FBigdecimal) string() string {
	return f.v.String()
}

type FBigint struct {
	v *big.Int
}

func newFBigint(v *big.Int, opt map[string]interface{}) *FBigint {
	return &FBigint{v: v}
}

func (f *FBigint) eql(v string) (bool, error) {
	value, ok := new(big.Int).SetString(v, 10)
	if !ok {
		return false, fmt.Errorf("failed to convert %q to bigint", v)
	}
	return f.v.Cmp(value) == 0, nil
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
			if sanitizeHex, found := opt["sanitize_hex"]; found {
				if sanitizeHex.(bool) {
					v = eth.SanitizeHex(v)
				}
			}
			return v
		},
	}
}

func (f *FString) eql(v string) (bool, error) {
	return f.v == f.process(v), nil
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

func (f *FBool) eql(v string) (bool, error) {
	value := v == "true"
	return value == f.v, nil
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

func (f *FBytes) eql(v string) (bool, error) {
	var data []byte
	var err error
	if strings.HasPrefix(v, "0x") {
		if data, err = hex.DecodeString(eth.SanitizeHex(v)); err != nil {
			return false, fmt.Errorf("failed to convert %s to byte array: %w", v, err)
		}
	}
	return bytes.Compare(f.v, data) == 0, nil
}

func (f *FBytes) string() string {
	return hex.EncodeToString(f.v)
}

type FArray struct {
	v []Comparable
}

func newFArray(v []Comparable) *FArray {
	return &FArray{v: v}
}

func (f *FArray) eql(v string) (bool, error) {
	in := newFArrayFromStr(v)
	if len(in) != len(f.v) {
		return false, nil
	}
	for i := 0; i < len(in); i++ {
		ok, err := f.v[i].eql(in[i].string())
		if err != nil {
			return false, fmt.Errorf("failed to compare elem %d: %w", err, i)
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

func newFArrayFromStr(v string) (out []Comparable) {
	if v == "" {
		return []Comparable{}
	}
	fmt.Println(">", v)
	panic(fmt.Errorf("need to implement this"))
}

func (f *FArray) string() string {
	strs := []string{}
	for _, field := range f.v {
		strs = append(strs, field.string())
	}
	return strings.Join(strs, ", ")
}
