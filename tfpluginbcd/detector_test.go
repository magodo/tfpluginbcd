package tfpluginbcd

import (
	"testing"

	"github.com/magodo/tfpluginschema/schema"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestCompareAttribute(t *testing.T) {
	cases := []struct {
		name   string
		scope  Scope
		path   []string
		oattr  *schema.Attribute
		nattr  *schema.Attribute
		expect []Change
	}{
		{
			name:  "Attribute deleted",
			scope: ResourceScope{Type: "foo_resource"},
			path:  []string{"attr1"},
			oattr: &schema.Attribute{
				Type:     cty.Bool,
				Required: true,
			},
			nattr: nil,
			expect: []Change{
				AttributeChange{
					Scope:    ResourceScope{Type: "foo_resource"},
					Path:     []string{"attr1"},
					IsDelete: true,
				},
			},
		},
		{
			name:  "Attribute added",
			scope: ResourceScope{Type: "foo_resource"},
			path:  []string{"attr1"},
			oattr: nil,
			nattr: &schema.Attribute{
				Type:     cty.Bool,
				Required: true,
			},
			expect: []Change{
				AttributeChange{
					Scope: ResourceScope{Type: "foo_resource"},
					Path:  []string{"attr1"},
					IsAdd: true,
					Current: &Attribute{
						Type:     cty.Bool,
						Required: true,
					},
				},
			},
		},
		{
			name:  "Attribute single updated",
			scope: ResourceScope{Type: "foo_resource"},
			path:  []string{"attr1"},
			oattr: &schema.Attribute{
				Type:     cty.Bool,
				Required: true,
			},
			nattr: &schema.Attribute{
				Type:     cty.String,
				Required: true,
			},
			expect: []Change{
				AttributeChange{
					Scope:    ResourceScope{Type: "foo_resource"},
					Path:     []string{"attr1"},
					IsModify: true,
					Current: &Attribute{
						Type:     cty.String,
						Required: true,
					},
					Modification: &AttributeModify{
						Type: &Modification[cty.Type]{
							From: cty.Bool,
							To:   cty.String,
						},
					},
				},
			},
		},
		{
			name:  "Attribute multiple updates",
			scope: ResourceScope{Type: "foo_resource"},
			path:  []string{"attr1"},
			oattr: &schema.Attribute{
				Type:     cty.Bool,
				Required: true,
			},
			nattr: &schema.Attribute{
				Type:     cty.String,
				Optional: true,
			},
			expect: []Change{
				AttributeChange{
					Scope:    ResourceScope{Type: "foo_resource"},
					Path:     []string{"attr1"},
					IsModify: true,
					Current: &Attribute{
						Type:     cty.String,
						Optional: true,
					},
					Modification: &AttributeModify{
						Type: &Modification[cty.Type]{
							From: cty.Bool,
							To:   cty.String,
						},
						Required: &Modification[bool]{
							From: true,
							To:   false,
						},
						Optional: &Modification[bool]{
							From: false,
							To:   true,
						},
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, compareAttribute(tt.scope, tt.path, tt.oattr, tt.nattr))
		})
	}
}

func TestCompareNestedBlock(t *testing.T) {
	cases := []struct {
		name   string
		scope  Scope
		path   []string
		oblk   *schema.NestedBlock
		nblk   *schema.NestedBlock
		expect []Change
	}{
		{
			name:  "NestedBlock deleted",
			scope: ResourceScope{Type: "foo_resource"},
			path:  []string{"blk1"},
			oblk: &schema.NestedBlock{
				NestingMode: schema.NestingSingle,
				Required:    true,
				Block: &schema.Block{
					Attributes:   map[string]*schema.Attribute{},
					NestedBlocks: map[string]*schema.NestedBlock{},
				},
			},
			nblk: nil,
			expect: []Change{
				BlockChange{
					Scope:    ResourceScope{Type: "foo_resource"},
					Path:     []string{"blk1"},
					IsDelete: true,
				},
			},
		},
		{
			name:  "NestedBlock added",
			scope: ResourceScope{Type: "foo_resource"},
			path:  []string{"blk1"},
			oblk:  nil,
			nblk: &schema.NestedBlock{
				NestingMode: schema.NestingSingle,
				Required:    true,
				Block: &schema.Block{
					Attributes:   map[string]*schema.Attribute{},
					NestedBlocks: map[string]*schema.NestedBlock{},
				},
			},
			expect: []Change{
				BlockChange{
					Scope: ResourceScope{Type: "foo_resource"},
					Path:  []string{"blk1"},
					IsAdd: true,
					Current: &Block{
						NestingMode: schema.NestingSingle,
						Required:    true,
					},
				},
			},
		},
		{
			name:  "NestedBlock single updated",
			scope: ResourceScope{Type: "foo_resource"},
			path:  []string{"blk1"},
			oblk: &schema.NestedBlock{
				NestingMode: schema.NestingSingle,
				Required:    true,
				Block: &schema.Block{
					Attributes:   map[string]*schema.Attribute{},
					NestedBlocks: map[string]*schema.NestedBlock{},
				},
			},
			nblk: &schema.NestedBlock{
				NestingMode: schema.NestingGroup,
				Required:    true,
				Block: &schema.Block{
					Attributes:   map[string]*schema.Attribute{},
					NestedBlocks: map[string]*schema.NestedBlock{},
				},
			},
			expect: []Change{
				BlockChange{
					Scope:    ResourceScope{Type: "foo_resource"},
					Path:     []string{"blk1"},
					IsModify: true,
					Current: &Block{
						NestingMode: schema.NestingGroup,
						Required:    true,
					},
					Modification: &BlockModify{
						NestingMode: &Modification[schema.NestingMode]{
							From: schema.NestingSingle,
							To:   schema.NestingGroup,
						},
					},
				},
			},
		},
		{
			name:  "NestedBlock multiple updates",
			scope: ResourceScope{Type: "foo_resource"},
			path:  []string{"blk1"},
			oblk: &schema.NestedBlock{
				NestingMode: schema.NestingSingle,
				Required:    true,
				Block: &schema.Block{
					Attributes:   map[string]*schema.Attribute{},
					NestedBlocks: map[string]*schema.NestedBlock{},
				},
			},
			nblk: &schema.NestedBlock{
				NestingMode: schema.NestingGroup,
				Optional:    true,
				Block: &schema.Block{
					Attributes:   map[string]*schema.Attribute{},
					NestedBlocks: map[string]*schema.NestedBlock{},
				},
			},
			expect: []Change{
				BlockChange{
					Scope:    ResourceScope{Type: "foo_resource"},
					Path:     []string{"blk1"},
					IsModify: true,
					Current: &Block{
						NestingMode: schema.NestingGroup,
						Optional:    true,
					},
					Modification: &BlockModify{
						NestingMode: &Modification[schema.NestingMode]{
							From: schema.NestingSingle,
							To:   schema.NestingGroup,
						},
						Required: &Modification[bool]{
							From: true,
							To:   false,
						},
						Optional: &Modification[bool]{
							From: false,
							To:   true,
						},
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, compareNestedBlock(tt.scope, tt.path, tt.oblk, tt.nblk))
		})
	}
}

func TestCompareBlock(t *testing.T) {
	cases := []struct {
		name   string
		scope  Scope
		path   []string
		oblk   *schema.Block
		nblk   *schema.Block
		expect []Change
	}{
		{
			name:  "Add/Delete/Modify attributes and nested blocks",
			scope: ResourceScope{Type: "foo_resource"},
			path:  []string{"blk"},
			oblk: &schema.Block{
				Attributes: map[string]*schema.Attribute{
					"old_only_attr": {
						Type: cty.Bool,
					},
					"update_attr": {
						Type: cty.Bool,
					},
				},
				NestedBlocks: map[string]*schema.NestedBlock{
					"old_only_blk": {
						NestingMode: schema.NestingSingle,
						Block: &schema.Block{
							Attributes:   map[string]*schema.Attribute{},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
					"update_blk": {
						NestingMode: schema.NestingSingle,
						Block: &schema.Block{
							Attributes:   map[string]*schema.Attribute{},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
				},
			},
			nblk: &schema.Block{
				Attributes: map[string]*schema.Attribute{
					"new_only_attr": {
						Type: cty.Bool,
					},
					"update_attr": {
						Type: cty.String,
					},
				},
				NestedBlocks: map[string]*schema.NestedBlock{
					"new_only_blk": {
						NestingMode: schema.NestingSingle,
						Block: &schema.Block{
							Attributes:   map[string]*schema.Attribute{},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
					"update_blk": {
						NestingMode: schema.NestingGroup,
						Block: &schema.Block{
							Attributes:   map[string]*schema.Attribute{},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
				},
			},
			expect: []Change{
				AttributeChange{
					Scope:    ResourceScope{Type: "foo_resource"},
					Path:     []string{"blk", "old_only_attr"},
					IsDelete: true,
				},
				AttributeChange{
					Scope:    ResourceScope{Type: "foo_resource"},
					Path:     []string{"blk", "update_attr"},
					IsModify: true,
					Current: &Attribute{
						Type: cty.String,
					},
					Modification: &AttributeModify{
						Type: &Modification[cty.Type]{
							From: cty.Bool,
							To:   cty.String,
						},
					},
				},
				AttributeChange{
					Scope: ResourceScope{Type: "foo_resource"},
					Path:  []string{"blk", "new_only_attr"},
					IsAdd: true,
					Current: &Attribute{
						Type: cty.Bool,
					},
				},
				BlockChange{
					Scope:    ResourceScope{Type: "foo_resource"},
					Path:     []string{"blk", "old_only_blk"},
					IsDelete: true,
				},
				BlockChange{
					Scope:    ResourceScope{Type: "foo_resource"},
					Path:     []string{"blk", "update_blk"},
					IsModify: true,
					Current: &Block{
						NestingMode: schema.NestingGroup,
					},
					Modification: &BlockModify{
						NestingMode: &Modification[schema.NestingMode]{
							From: schema.NestingSingle,
							To:   schema.NestingGroup,
						},
					},
				},
				BlockChange{
					Scope: ResourceScope{Type: "foo_resource"},
					Path:  []string{"blk", "new_only_blk"},
					IsAdd: true,
					Current: &Block{
						NestingMode: schema.NestingSingle,
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, compareBlock(tt.scope, tt.path, tt.oblk, tt.nblk))
		})
	}
}

func TestCompareResources(t *testing.T) {
	cases := []struct {
		name   string
		orm    map[string]*schema.Resource
		nrm    map[string]*schema.Resource
		expect []Change
	}{
		{
			name: "Resource deleted",
			orm: map[string]*schema.Resource{
				"foo_resource": {
					Block: &schema.Block{
						Attributes:   map[string]*schema.Attribute{},
						NestedBlocks: map[string]*schema.NestedBlock{},
					},
				},
			},
			nrm: map[string]*schema.Resource{},
			expect: []Change{
				ResourceChange{
					Type:     "foo_resource",
					IsDelete: true,
				},
			},
		},
		{
			name: "Resource added",
			orm:  map[string]*schema.Resource{},
			nrm: map[string]*schema.Resource{
				"foo_resource": {
					Block: &schema.Block{
						Attributes:   map[string]*schema.Attribute{},
						NestedBlocks: map[string]*schema.NestedBlock{},
					},
				},
			},
			expect: []Change{
				ResourceChange{
					Type:    "foo_resource",
					IsAdd:   true,
					Current: &Resource{},
				},
			},
		},
		{
			name: "Resource update",
			orm: map[string]*schema.Resource{
				"foo_resource": {
					Block: &schema.Block{
						Attributes:   map[string]*schema.Attribute{},
						NestedBlocks: map[string]*schema.NestedBlock{},
					},
				},
			},
			nrm: map[string]*schema.Resource{
				"foo_resource": {
					SchemaVersion: 1,
					Block: &schema.Block{
						Attributes:   map[string]*schema.Attribute{},
						NestedBlocks: map[string]*schema.NestedBlock{},
					},
				},
			},
			expect: []Change{
				ResourceChange{
					Type:     "foo_resource",
					IsModify: true,
					Current: &Resource{
						SchemaVersion: 1,
					},
					Modification: &ResourceModify{
						SchemaVersion: &Modification[int]{
							From: 0,
							To:   1,
						},
					},
				},
			},
		},
		{
			name: "Resource internal (attr/block) update",
			orm: map[string]*schema.Resource{
				"foo_resource": {
					Block: &schema.Block{
						Attributes: map[string]*schema.Attribute{
							"old_only_attr": {},
						},
						NestedBlocks: map[string]*schema.NestedBlock{},
					},
				},
			},
			nrm: map[string]*schema.Resource{
				"foo_resource": {
					Block: &schema.Block{
						Attributes:   map[string]*schema.Attribute{},
						NestedBlocks: map[string]*schema.NestedBlock{},
					},
				},
			},
			expect: []Change{
				AttributeChange{
					Scope:    ResourceScope{Type: "foo_resource"},
					Path:     []string{"old_only_attr"},
					IsDelete: true,
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, compareResources(tt.orm, tt.nrm, false))
		})
	}
}

func TestCompare(t *testing.T) {
	cases := []struct {
		name   string
		osch   *schema.ProviderSchema
		nsch   *schema.ProviderSchema
		expect []Change
	}{
		{
			name: "Provider config deleted (though not gonna happen in real life)",
			osch: &schema.ProviderSchema{
				Provider: &schema.Schema{
					Block: &schema.Block{
						Attributes:   map[string]*schema.Attribute{},
						NestedBlocks: map[string]*schema.NestedBlock{},
					},
				},
			},
			nsch: &schema.ProviderSchema{},
			expect: []Change{
				ProviderChange{
					IsDelete: true,
				},
			},
		},
		{
			name: "Provider config added (though not gonna happen in real life)",
			osch: &schema.ProviderSchema{},
			nsch: &schema.ProviderSchema{
				Provider: &schema.Schema{
					Block: &schema.Block{
						Attributes:   map[string]*schema.Attribute{},
						NestedBlocks: map[string]*schema.NestedBlock{},
					},
				},
			},
			expect: []Change{
				ProviderChange{
					IsAdd: true,
				},
			},
		},
		{
			name: "Provider config attribute deleted",
			osch: &schema.ProviderSchema{
				Provider: &schema.Schema{
					Block: &schema.Block{
						Attributes: map[string]*schema.Attribute{
							"old_only_attr": {},
						},
						NestedBlocks: map[string]*schema.NestedBlock{},
					},
				},
			},
			nsch: &schema.ProviderSchema{
				Provider: &schema.Schema{
					Block: &schema.Block{
						Attributes:   map[string]*schema.Attribute{},
						NestedBlocks: map[string]*schema.NestedBlock{},
					},
				},
			},
			expect: []Change{
				AttributeChange{
					Scope:    ProviderScope{},
					Path:     []string{"old_only_attr"},
					IsDelete: true,
				},
			},
		},
		{
			name: "Resource/DataSource attribute deleted",
			osch: &schema.ProviderSchema{
				DataSourceSchemas: map[string]*schema.Resource{
					"foo_resource": {
						Block: &schema.Block{
							Attributes: map[string]*schema.Attribute{
								"old_only_attr": {},
							},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
				},
				ResourceSchemas: map[string]*schema.Resource{
					"bar_resource": {
						Block: &schema.Block{
							Attributes: map[string]*schema.Attribute{
								"old_only_attr": {},
							},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
				},
			},
			nsch: &schema.ProviderSchema{
				DataSourceSchemas: map[string]*schema.Resource{
					"foo_resource": {
						Block: &schema.Block{
							Attributes:   map[string]*schema.Attribute{},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
				},
				ResourceSchemas: map[string]*schema.Resource{
					"bar_resource": {
						Block: &schema.Block{
							Attributes:   map[string]*schema.Attribute{},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					},
				},
			},
			expect: []Change{
				AttributeChange{
					Scope:    DataSourceScope{Type: "foo_resource"},
					Path:     []string{"old_only_attr"},
					IsDelete: true,
				},
				AttributeChange{
					Scope:    ResourceScope{Type: "bar_resource"},
					Path:     []string{"old_only_attr"},
					IsDelete: true,
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, Compare(tt.osch, tt.nsch))
		})
	}
}
