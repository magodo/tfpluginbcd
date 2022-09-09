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
		rule    string
		expect  []Change
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
			rule: `
c.kind == "resource"
c.is_add
`,
			expect: []Change{
				ResourceChange{
					Type:         "foo_resource",
					IsDataSource: false,
					IsAdd:        true,
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := Filter(context.TODO(), tt.changes, buildRegoModule(buildRule(tt.rule)))
			require.NoError(t, err)
			require.Equal(t, tt.expect, actual)
		})
	}
}
