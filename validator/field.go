package validator

import (
	"fmt"

	pbentities "github.com/streamingfast/substreams-test/pb/entity/v1"
	"github.com/streamingfast/substreams-test/validator/fields"
)

func (v *Validator) newField(substreamsEntity string, field *pbentities.Field) *fields.Field {
	graphqlFieldName := v.getGraphQLFieldName(substreamsEntity, field.Name)
	out := &fields.Field{
		SubstreamsField: field.Name,
		GraphqlField:    graphqlFieldName,
		GraphqlJSONPath: fmt.Sprintf("data.%s.%s", normalizeEntityName(substreamsEntity), graphqlFieldName),
	}

	if v.isGraphQLAssociatedField(substreamsEntity, field.Name) {
		out.GraphqlField = fmt.Sprintf("%s { id }", graphqlFieldName)
		if v.isGraphQLArrayField(substreamsEntity, field.Name) {
			out.GraphqlJSONPath = fmt.Sprintf("data.%s.%s.#.id", normalizeEntityName(substreamsEntity), graphqlFieldName)
		} else {
			out.GraphqlJSONPath = fmt.Sprintf("data.%s.%s.id", normalizeEntityName(substreamsEntity), graphqlFieldName)
		}
	}

	fieldOpt := v.getFieldOpt(substreamsEntity, field.Name)
	out.Obj, out.ObjFactory = fields.ParseValue(field.NewValue, fieldOpt)

	return out
}

func normalizeEntityName(s string) string {
	if len(s) != 0 && (s[0] <= 90 && s[0] >= 65) {
		return string(s[0]+32) + string(s[1:])
	}
	return s
}
