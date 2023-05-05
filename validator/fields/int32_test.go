package fields

import (
	pbentities "github.com/streamingfast/substreams-test/pb/entity/v1"
	"github.com/streamingfast/substreams-test/validator/config"
	"github.com/test-go/testify/assert"
	"testing"
)

func Test_int32(t *testing.T) {
	tests := []struct {
		name         string
		entityFields []*Field
		expected     bool
	}{
		{
			name: "int32 value equal",
			entityFields: []*Field{
				NewField(&pbentities.Value{
					Typed: &pbentities.Value_Int32{Int32: 10},
				}, "10"),
			},
			expected: true,
		},
		{
			name: "int32 value not equal",
			entityFields: []*Field{
				NewField(&pbentities.Value{
					Typed: &pbentities.Value_Int32{Int32: 10},
				}, "20"),
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, entityField := range test.entityFields {
				graphValue, _ := entityField.ObjFactory(entityField.graphEntity, config.Options{})
				assert.Equal(t, entityField.Obj.Eql(graphValue), test.expected)
			}
		})
	}
}
