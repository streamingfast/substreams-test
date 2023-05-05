package fields

import (
	pbentities "github.com/streamingfast/substreams-test/pb/entity/v1"
	"github.com/streamingfast/substreams-test/validator/config"
	"github.com/test-go/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		name         string
		entityFields []*Field
		expected     bool
		expectedErr  bool
	}{
		{
			name: "equal string arrays in the same order",
			entityFields: []*Field{
				NewField(&pbentities.Value{
					Typed: &pbentities.Value_Array{
						Array: &pbentities.Array{
							Value: []*pbentities.Value{
								{
									Typed: &pbentities.Value_String_{
										String_: "0xaa",
									},
								},
								{
									Typed: &pbentities.Value_String_{
										String_: "0xbb",
									},
								},
								{
									Typed: &pbentities.Value_String_{
										String_: "0xcc",
									},
								},
							},
						},
					},
				}, "[\"0xaa\", \"0xbb\", \"0xcc\"]"),
			},
			expected: true,
		},
		{
			name: "equal string arrays in the different order",
			entityFields: []*Field{
				NewField(&pbentities.Value{
					Typed: &pbentities.Value_Array{
						Array: &pbentities.Array{
							Value: []*pbentities.Value{
								{
									Typed: &pbentities.Value_String_{
										String_: "0xaa",
									},
								},
								{
									Typed: &pbentities.Value_String_{
										String_: "0xbb",
									},
								},
								{
									Typed: &pbentities.Value_String_{
										String_: "0xcc",
									},
								},
							},
						},
					},
				}, "[\"0xbb\", \"0xaa\", \"0xcc\"]"),
			},
			expected: true,
		},
		{
			name: "not equal length string arrays",
			entityFields: []*Field{
				NewField(&pbentities.Value{
					Typed: &pbentities.Value_Array{
						Array: &pbentities.Array{
							Value: []*pbentities.Value{
								{
									Typed: &pbentities.Value_String_{
										String_: "0xaa",
									},
								},
								{
									Typed: &pbentities.Value_String_{
										String_: "0xbb",
									},
								},
								{
									Typed: &pbentities.Value_String_{
										String_: "0xcc",
									},
								},
							},
						},
					},
				}, "[\"0xbb\", \"0xcc\"]"),
			},
			expectedErr: true,
		},
		{
			name: "equal int arrays",
			entityFields: []*Field{
				NewField(&pbentities.Value{
					Typed: &pbentities.Value_Array{
						Array: &pbentities.Array{
							Value: []*pbentities.Value{
								{
									Typed: &pbentities.Value_Int32{
										Int32: 10,
									},
								},
								{
									Typed: &pbentities.Value_Int32{
										Int32: 20,
									},
								},
								{
									Typed: &pbentities.Value_Int32{
										Int32: 30,
									},
								},
							},
						},
					},
				}, "[10, 20, 30]"),
			},
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, entityField := range test.entityFields {
				graphValue, err := entityField.ObjFactory(entityField.graphEntity, config.Options{})
				if test.expectedErr {
					assert.Error(t, err)
				} else {
					assert.Equal(t, entityField.Obj.Eql(graphValue), test.expected)
				}
			}
		})
	}
}
