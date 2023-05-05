package fields

import (
	pbentities "github.com/streamingfast/substreams-test/pb/entity/v1"
	"github.com/streamingfast/substreams-test/validator/config"
	"github.com/test-go/testify/assert"
	"math/big"
	"testing"
)

func Test_decimal_eql(t *testing.T) {
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
				graphValue, _ := entityField.ObjFactory(entityField.graphEntity, config.Options{})
				assert.Equal(t, entityField.Obj.Eql(graphValue), test.expected)
			}
		})
	}
}

func Test_decimal_options(t *testing.T) {
	ignoreError := float64(-1)
	ignoreTolerance := float64(-1)
	ignoreRound := ""

	tests := []struct {
		name          string
		entityField   *Field
		bigDecimal    *Decimal
		expectedEqual bool
	}{
		{
			name: "no config",
			entityField: NewField(&pbentities.Value{
				Typed: &pbentities.Value_Bigdecimal{Bigdecimal: "3527.243688533782218063279378930382749686619384024845485673704637885598213415263657383626782854522193"},
			}, "3527.243688533782350769224945899756"),
			bigDecimal: &Decimal{
				Options: config.Options{},
				v:       new(big.Float).SetFloat64(3527.243688533782218063279378930382749686619384024845485673704637885598213415263657383626782854522193),
				str:     "3527.243688533782218063279378930382749686619384024845485673704637885598213415263657383626782854522193",
			},
			expectedEqual: false,
		},
		{
			name: "error at 0.0001",
			entityField: NewField(&pbentities.Value{
				Typed: &pbentities.Value_Bigdecimal{Bigdecimal: "3527.243688533782218063279378930382749686619384024845485673704637885598213415263657383626782854522193"},
			}, "3527.243688533782350769224945899756"),
			bigDecimal: &Decimal{
				Options: config.NewOptions(0.0001, ignoreTolerance, ignoreRound),
				v:       new(big.Float).SetFloat64(3527.243688533782218063279378930382749686619384024845485673704637885598213415263657383626782854522193),
				str:     "3527.243688533782218063279378930382749686619384024845485673704637885598213415263657383626782854522193",
			},
			expectedEqual: true,
		},
		{
			name: "round shortest",
			entityField: NewField(&pbentities.Value{
				Typed: &pbentities.Value_Bigdecimal{Bigdecimal: "3527.011111"},
			}, "3527.01"),
			bigDecimal: &Decimal{
				Options: config.NewOptions(ignoreError, ignoreTolerance, "shortest"),
				v:       new(big.Float).SetFloat64(3527.243688533782218063279378930382749686619384024845485673704637885598213415263657383626782854522193),
				str:     "3527.011111",
			},
			expectedEqual: true,
		},
		{
			name: "tolerance",
			entityField: NewField(&pbentities.Value{
				Typed: &pbentities.Value_Bigdecimal{Bigdecimal: "3371.050771830142198236003792915391515091655972474595213881058766601526391568084628427376630197053540"},
			}, "3371.050771830142325065483050114257"),
			bigDecimal: &Decimal{
				Options: config.NewOptions(0.001, ignoreTolerance, ignoreRound),
				v:       new(big.Float).SetFloat64(3371.050771830142198236003792915391515091655972474595213881058766601526391568084628427376630197053540),
				str:     "3371.050771830142198236003792915391515091655972474595213881058766601526391568084628427376630197053540",
			},
			expectedEqual: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entityField := test.entityField
			graphValue, _ := entityField.ObjFactory(test.entityField.graphEntity, config.Options{})
			bigDecimal := test.bigDecimal
			assert.Equal(t, bigDecimal.Eql(graphValue), test.expectedEqual)
		})
	}
}
