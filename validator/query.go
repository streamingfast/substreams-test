package validator

import (
	"fmt"
	"strings"
)

func queryFromEntity(entity string, fields []*Field) string {
	sfield := []string{}
	for _, f := range fields {
		sfield = append(sfield, f.graphqlField)
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
