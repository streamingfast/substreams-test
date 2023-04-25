package fields

import (
	b64 "encoding/base64"
	"fmt"
	"math/big"

	pbentities "github.com/streamingfast/substreams-test/pb/entity/v1"
)

type Field struct {
	substreamsEntity string
	SubstreamsField  string
	graphEntity      string
	GraphqlField     string
	GraphqlJSONPath  string

	ObjFactory ComparableFactory
	Obj        Comparable
}

type ComparableFactory func(string) (Comparable, error)
type Comparable interface {
	Eql(v Comparable) bool
	String() string
}

func ParseValue(value *pbentities.Value, fieldOpt map[string]interface{}) (Comparable, ComparableFactory) {
	switch newValue := value.Typed.(type) {
	case *pbentities.Value_Int32:
		return newInt32(newValue.Int32, fieldOpt), newInt32FromStr

	case *pbentities.Value_Bigdecimal:
		nvalue, _ := new(big.Float).SetString(newValue.Bigdecimal)
		return newDecimal(nvalue, fieldOpt), newDecimalFromStr

	case *pbentities.Value_Bigint:
		nvalue, _ := new(big.Int).SetString(newValue.Bigint, 10)
		return newBigint(nvalue, fieldOpt), newBigintFromStr

	case *pbentities.Value_String_:
		return newString(newValue.String_, fieldOpt), newStringFromStr

	case *pbentities.Value_Bool:
		return newBool(newValue.Bool, fieldOpt), newFBoolFromStr

	case *pbentities.Value_Bytes:
		data, err := b64.StdEncoding.DecodeString(newValue.Bytes)
		if err != nil {
			panic(err)
		}

		return newBytes(data, fieldOpt), newBytesFromStr
	case *pbentities.Value_Array:
		arr := &Array{}
		arrFactory := &FArrayFactory{}
		for _, v := range newValue.Array.Value {
			val, factory := ParseValue(v, fieldOpt)
			arr.values = append(arr.values, val)
			arrFactory.factories = append(arrFactory.factories, factory)
		}
		return arr, arrFactory.newArrayFromStr
	default:
		panic(fmt.Errorf("unknown field v value type %T", newValue))
	}

}
