package subcmp

import (
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
	"testing"
)

func Test_extractGraphqlVariables(t *testing.T) {
	tests := []struct {
		query       string
		expect      []string
		expectError bool
	}{
		{
			query:  `graphqlQuery($block: Int!,$pool: String!)  { pools(where: {id: $pool},block: {number: $block}) { id } }`,
			expect: []string{"block", "pool"},
		},
		{
			query:  `graphqlQuery()  { pools(where: {id: $pool},block: {number: $block}) { id } }`,
			expect: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.query, func(t *testing.T) {
			out, err := extractGraphqlVariables(test.query)
			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expect, out)
			}
		})
	}
}

func Test_extractSubstreamsPathVariables(t *testing.T) {
	tests := []struct {
		path   string
		expect []string
	}{
		{
			path:   ".feeGrowthGlobalUpdates[] | select(.poolAddress == \"${pool}\") | .newValue.value",
			expect: []string{"pool"},
		},
		{
			path:   ".feeGrowthGlobalUpdates[] | select(.poolAddress == \"${pool}\") | .${foo}.value",
			expect: []string{"pool", "foo"},
		},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			assert.Equal(t, test.expect, extractSubstreamsPathVariables(test.path))
		})
	}
}

func Test_replaceVarsSubstreamPath(t *testing.T) {
	tests := []struct {
		vars          map[string]string
		substreamPath string
		expect        string
	}{
		{
			vars: map[string]string{
				"pool": "aa",
			},
			substreamPath: ".feeGrowthGlobalUpdates[] | select(.poolAddress == \"${pool}\") | .newValue.value",
			expect:        ".feeGrowthGlobalUpdates[] | select(.poolAddress == \"aa\") | .newValue.value",
		},
		{
			vars: map[string]string{
				"pool": "aa",
				"foo":  "bb",
			},
			substreamPath: ".feeGrowthGlobalUpdates[] | select(.poolAddress == \"${pool}\") | .${foo}.value",
			expect:        ".feeGrowthGlobalUpdates[] | select(.poolAddress == \"aa\") | .bb.value",
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, test.expect, replaceVarsSubstreamPath(test.vars, test.substreamPath))
		})
	}
}
