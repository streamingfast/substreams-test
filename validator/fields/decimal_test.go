package fields

import (
	pbentities "github.com/streamingfast/substreams-test/pb/entity/v1"
	"github.com/test-go/testify/assert"
	"testing"
)

func Test_decimal(t *testing.T) {
	tests := []struct {
		name         string
		entityFields []*Field
		expected     bool
	}{
		{
			name: "decimal value equal",
			entityFields: []*Field{
				NewField(&pbentities.Value{
					Typed: &pbentities.Value_Bigdecimal{Bigdecimal: "923764728342093848726417618297313901273019792136823410237412976834"},
				}, "923764728342093848726417618297313901273019792136823410237412976834"),
			},
			expected: true,
		},
		{
			name: "decimal value not equal",
			entityFields: []*Field{
				NewField(&pbentities.Value{
					Typed: &pbentities.Value_Bigdecimal{Bigdecimal: "923764728342093848726417618297313901273019792136823410237412976834"},
				}, "1"),
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, entityField := range test.entityFields {
				graphValue, _ := entityField.ObjFactory(entityField.graphEntity)
				assert.Equal(t, entityField.Obj.Eql(graphValue), test.expected)
			}
		})
	}
}
