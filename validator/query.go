package validator

import (
	"fmt"
	"strings"

	"github.com/streamingfast/substreams-test/validator/fields"
)

func queryFromEntity(entity string, fields []*fields.Field) string {
	sfield := []string{}
	for _, f := range fields {
		sfield = append(sfield, f.GraphqlField)
	}
	return fmt.Sprintf(`
query($block: Int!,$id: String!) {
	%s(id: $id,block: {number: $block}) {
		id
		%s
	}
}
`, entity, strings.Join(sfield, "\n"))
}
