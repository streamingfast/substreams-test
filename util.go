package subcmp

import (
	"fmt"
	"regexp"
	"strings"
)

var graphqlQueryRegEx = regexp.MustCompile("^query\\((.*?\\))")
var graphqlVarRegEx = regexp.MustCompile("(\\$[a-zA-Z]*)")

func extractGraphqlVariables(query string) ([]string, error) {
	matches := graphqlQueryRegEx.FindAllString(query, -1)
	if len(matches) != 1 {
		return nil, fmt.Errorf("failed to extract graphqlQuery function")
	}

	var out []string
	vars := graphqlVarRegEx.FindAllString(matches[0], -1)
	for _, match := range vars {
		out = append(out, strings.Replace(match, "$", "", -1))
	}
	return out, nil
}

var substreamsVarRegEx = regexp.MustCompile("\\$\\{([a-zA-Z]*\\})")

func extractSubstreamsPathVariables(substreamPath string) (out []string) {
	varNames := substreamsVarRegEx.FindAllString(substreamPath, -1)
	for _, varName := range varNames {
		out = append(out, strings.Replace(strings.Replace(varName, "${", "", -1), "}", "", -1))
	}
	return out
}

func replaceVarsSubstreamPath(vars map[string]string, substreamsPath string) string {
	path := substreamsPath
	for varName, varValue := range vars {
		path = strings.Replace(path, fmt.Sprintf("${%s}", varName), varValue, -1)
	}
	return path
}
