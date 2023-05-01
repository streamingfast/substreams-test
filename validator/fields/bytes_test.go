package fields

import (
	pbentities "github.com/streamingfast/substreams-test/pb/entity/v1"
	"github.com/test-go/testify/assert"
	"testing"
)

func Test_bytes(t *testing.T) {
	tests := []struct {
		name         string
		entityFields []*Field
		expected     bool
	}{
		{
			name: "bytes value equal",
			entityFields: []*Field{
				NewField(&pbentities.Value{
					Typed: &pbentities.Value_Bytes{Bytes: "qg=="},
				}, "0xaa"),
			},
			expected: true,
		},
		{
			name: "bytes value not equal",
			entityFields: []*Field{
				NewField(&pbentities.Value{
					Typed: &pbentities.Value_Bytes{Bytes: "qg=="},
				}, "0xbb"),
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
