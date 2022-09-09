package tfpluginbcd

import (
	"context"
	"testing"

	"github.com/magodo/tfpluginschema/schema"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestRun(t *testing.T) {
	cases := []struct {
		name       string
		opt        Opt
		osch, nsch schema.ProviderSchema
		filtN      int
		hasError   bool
	}{
		{
			name: "not defined rule",
			opt: Opt{
				Rules: []string{"Rxxx"},
			},
			hasError: true,
		},
		{
			name: "rule1",
			opt: Opt{
				Rules: []string{"R001"},
			},
			osch: schema.ProviderSchema{
				ResourceSchemas: map[string]*schema.Resource{
					"foo_resource": {},
				},
			},
			nsch:  schema.ProviderSchema{},
			filtN: 1,
		},
		{
			name: "rule2",
			opt: Opt{
				Rules: []string{"R002"},
			},
			osch: schema.ProviderSchema{
				DataSourceSchemas: map[string]*schema.Resource{
					"foo_resource": {},
				},
			},
			nsch:  schema.ProviderSchema{},
			filtN: 1,
		},
		{
			name: "rule3",
			opt: Opt{
				Rules: []string{"R003"},
			},
			osch: schema.ProviderSchema{
				ResourceSchemas: map[string]*schema.Resource{
					"foo_resource": {
						Block: &schema.Block{
							Attributes: map[string]*schema.Attribute{
								"old_only_attr": {},
							},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
				},
			},
			nsch: schema.ProviderSchema{
				ResourceSchemas: map[string]*schema.Resource{
					"foo_resource": {
						Block: &schema.Block{
							Attributes:   map[string]*schema.Attribute{},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
				},
			},
			filtN: 1,
		},
		{
			name: "rule4",
			opt: Opt{
				Rules: []string{"R004"},
			},
			osch: schema.ProviderSchema{
				ResourceSchemas: map[string]*schema.Resource{
					"foo_resource": {
						Block: &schema.Block{
							Attributes: map[string]*schema.Attribute{},
							NestedBlocks: map[string]*schema.NestedBlock{
								"old_only_blk": {},
							},
						},
					},
				},
			},
			nsch: schema.ProviderSchema{
				ResourceSchemas: map[string]*schema.Resource{
					"foo_resource": {
						Block: &schema.Block{
							Attributes:   map[string]*schema.Attribute{},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
				},
			},
			filtN: 1,
		},
		{
			name: "rule5",
			opt: Opt{
				Rules: []string{"R005"},
			},
			osch: schema.ProviderSchema{
				ResourceSchemas: map[string]*schema.Resource{
					"foo_resource": {
						Block: &schema.Block{
							Attributes: map[string]*schema.Attribute{
								"attr": {
									Type: cty.Bool,
								},
							},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
				},
			},
			nsch: schema.ProviderSchema{
				ResourceSchemas: map[string]*schema.Resource{
					"foo_resource": {
						Block: &schema.Block{
							Attributes: map[string]*schema.Attribute{
								"attr": {
									Type: cty.String,
								},
							},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
				},
			},
			filtN: 1,
		},
		{
			name: "rule6",
			opt: Opt{
				Rules: []string{"R006"},
			},
			osch: schema.ProviderSchema{
				ResourceSchemas: map[string]*schema.Resource{
					"foo_resource": {
						Block: &schema.Block{
							Attributes: map[string]*schema.Attribute{
								"attr": {
									Type: cty.Bool,
								},
							},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
				},
			},
			nsch: schema.ProviderSchema{
				ResourceSchemas: map[string]*schema.Resource{
					"foo_resource": {
						Block: &schema.Block{
							Attributes: map[string]*schema.Attribute{
								"attr": {
									Type:     cty.Bool,
									Required: true,
								},
							},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
				},
			},
			filtN: 1,
		},
		{
			name: "rule7",
			opt: Opt{
				Rules: []string{"R007"},
			},
			osch: schema.ProviderSchema{
				ResourceSchemas: map[string]*schema.Resource{
					"foo_resource": {
						Block: &schema.Block{
							Attributes: map[string]*schema.Attribute{},
							NestedBlocks: map[string]*schema.NestedBlock{
								"blk": {
									Block: &schema.Block{
										Attributes:   map[string]*schema.Attribute{},
										NestedBlocks: map[string]*schema.NestedBlock{},
									},
								},
							},
						},
					},
				},
			},
			nsch: schema.ProviderSchema{
				ResourceSchemas: map[string]*schema.Resource{
					"foo_resource": {
						Block: &schema.Block{
							Attributes: map[string]*schema.Attribute{},
							NestedBlocks: map[string]*schema.NestedBlock{
								"blk": {
									Required: true,
									Block: &schema.Block{
										Attributes:   map[string]*schema.Attribute{},
										NestedBlocks: map[string]*schema.NestedBlock{},
									},
								},
							},
						},
					},
				},
			},
			filtN: 1,
		},
		{
			name: "rule8",
			opt: Opt{
				Rules: []string{"R008"},
			},
			osch: schema.ProviderSchema{
				ResourceSchemas: map[string]*schema.Resource{
					"foo_resource": {
						Block: &schema.Block{
							Attributes:   map[string]*schema.Attribute{},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
				},
			},
			nsch: schema.ProviderSchema{
				ResourceSchemas: map[string]*schema.Resource{
					"foo_resource": {
						Block: &schema.Block{
							Attributes: map[string]*schema.Attribute{
								"attr": {
									Type:     cty.Bool,
									Required: true,
								},
							},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
				},
			},
			filtN: 1,
		},
		{
			name: "rule9",
			opt: Opt{
				Rules: []string{"R009"},
			},
			osch: schema.ProviderSchema{
				ResourceSchemas: map[string]*schema.Resource{
					"foo_resource": {
						Block: &schema.Block{
							Attributes:   map[string]*schema.Attribute{},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
				},
			},
			nsch: schema.ProviderSchema{
				ResourceSchemas: map[string]*schema.Resource{
					"foo_resource": {
						Block: &schema.Block{
							Attributes: map[string]*schema.Attribute{},
							NestedBlocks: map[string]*schema.NestedBlock{
								"blk": {
									Required: true,
									Block: &schema.Block{
										Attributes:   map[string]*schema.Attribute{},
										NestedBlocks: map[string]*schema.NestedBlock{},
									},
								},
							},
						},
					},
				},
			},
			filtN: 1,
		},
		{
			name: "rule1 no match",
			opt: Opt{
				Rules: []string{"R001"},
			},
			osch: schema.ProviderSchema{
				ResourceSchemas: map[string]*schema.Resource{},
			},
			nsch: schema.ProviderSchema{
				ResourceSchemas: map[string]*schema.Resource{},
			},
			filtN: 0,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := run(context.TODO(), tt.osch, tt.nsch, tt.opt)
			if tt.hasError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.filtN, len(actual))
		})
	}
}
