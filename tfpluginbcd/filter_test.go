package tfpluginbcd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	cases := []struct {
		name    string
		changes []Change
		rules   []Rule
		expect  []FilterResult
	}{
		{
			name: "Resoruce is added",
			changes: []Change{
				ResourceChange{
					Type:         "foo_resource",
					IsDataSource: false,
					IsAdd:        true,
				},
				ResourceChange{
					Type:     "foo_resource",
					IsDelete: true,
				},
			},
			rules: []Rule{
				{
					ID: "CUSTOM",
					Expr: `
c.kind == "resource"
c.is_add
`,
				},
			},
			expect: []FilterResult{
				{
					Rule: "CUSTOM",
					Change: ResourceChange{
						Type:         "foo_resource",
						IsDataSource: false,
						IsAdd:        true,
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := Filter(context.TODO(), tt.changes, tt.rules)
			require.NoError(t, err)
			require.Equal(t, tt.expect, actual)
		})
	}
}
