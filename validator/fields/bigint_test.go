package fields

import (
	pbentities "github.com/streamingfast/substreams-test/pb/entity/v1"
	"github.com/streamingfast/substreams-test/validator/config"
	"github.com/test-go/testify/assert"
	"testing"
)

func Test_BigInt(t *testing.T) {
	tests := []struct {
		name         string
		entityFields []*Field
		expected     bool
	}{
		{
			name: "big int values equal",
			entityFields: []*Field{
				NewField(&pbentities.Value{
					Typed: &pbentities.Value_Bigint{
						Bigint: "10",
					},
				}, "10"),
			},
			expected: true,
		},
		{
			name: "big int values not equal",
			entityFields: []*Field{
				NewField(&pbentities.Value{
					Typed: &pbentities.Value_Bigint{
						Bigint: "10",
					},
				}, "15"),
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
