package tfpluginbcd

import (
	"testing"

	"github.com/magodo/tfpluginschema/schema"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestChangeString(t *testing.T) {
	cases := []struct {
		name   string
		change Change
		expect string
	}{
		{
			name: "Provider add",
			change: ProviderChange{
				IsAdd: true,
			},
			expect: "Provider config is added",
		},
		{
			name: "Provider delete",
			change: ProviderChange{
				IsDelete: true,
			},
			expect: "Provider config is deleted",
		},
		{
			name: "Resource add",
			change: ResourceChange{
				Type:         "foo_resource",
				IsDataSource: false,
				IsAdd:        true,
			},
			expect: "Resource foo_resource is added",
		},
		{
			name: "Data source add",
			change: ResourceChange{
				Type:         "foo_resource",
				IsDataSource: true,
				IsAdd:        true,
			},
			expect: "Data Source foo_resource is added",
		},
		{
			name: "Resource delete",
			change: ResourceChange{
				Type:     "foo_resource",
				IsDelete: true,
			},
			expect: "Resource foo_resource is deleted",
		},
		{
			name: "Resource change",
			change: ResourceChange{
				Type:     "foo_resource",
				IsModify: true,
				Modification: &ResourceModify{
					SchemaVersion: &Modification[int]{
						From: 0,
						To:   1,
					},
				},
			},
			expect: "Resource foo_resource is changed: schema version: 0 -> 1",
		},
		{
			name: "Provider config attribute add",
			change: AttributeChange{
				Scope: ProviderScope{},
				Path:  []string{"foo", "bar"},
				IsAdd: true,
			},
			expect: `Attribute "foo.bar" of provider config is added`,
		},
		{
			name: "Resource attribute add",
			change: AttributeChange{
				Scope: ResourceScope{Type: "foo_resource"},
				Path:  []string{"foo", "bar"},
				IsAdd: true,
			},
			expect: `Attribute "foo.bar" of resource foo_resource is added`,
		},
		{
			name: "Data Source attribute add",
			change: AttributeChange{
				Scope: DataSourceScope{Type: "foo_resource"},
				Path:  []string{"foo", "bar"},
				IsAdd: true,
			},
			expect: `Attribute "foo.bar" of data source foo_resource is added`,
		},
		{
			name: "Resource attribute delete",
			change: AttributeChange{
				Scope:    ResourceScope{Type: "foo_resource"},
				Path:     []string{"foo", "bar"},
				IsDelete: true,
			},
			expect: `Attribute "foo.bar" of resource foo_resource is deleted`,
		},
		{
			name: "Resource attribute change (complete)",
			change: AttributeChange{
				Scope:    ResourceScope{Type: "foo_resource"},
				Path:     []string{"foo", "bar"},
				IsModify: true,
				Modification: &AttributeModify{
					Type: &Modification[cty.Type]{
						From: cty.Bool,
						To:   cty.String,
					},
					Required: &Modification[bool]{
						From: false,
						To:   true,
					},
					Optional: &Modification[bool]{
						From: true,
						To:   false,
					},
					Computed: &Modification[bool]{
						From: false,
						To:   true,
					},
					// Intentionally ignore this
					// ForceNew: &Modification[bool]{
					// 	From: false,
					// 	To:   false,
					// },
					Default: &Modification[any]{
						From: nil,
						To:   10,
					},
					Sensitive: &Modification[bool]{
						From: false,
						To:   true,
					},
					ConflictsWith: &Modification[[]string]{
						From: []string{},
						To:   []string{"a"},
					},
					RequiredWith: &Modification[[]string]{
						From: []string{"a"},
						To:   []string{"b"},
					},
					ExactlyOneOf: &Modification[[]string]{
						From: nil,
						To:   []string{"a"},
					},
					AtLeastOneOf: &Modification[[]string]{
						From: []string{},
						To:   nil,
					},
				},
			},
			expect: `Attribute "foo.bar" of resource foo_resource is changed: ` +
				"type: bool -> string, " +
				"required: false -> true, " +
				"optional: true -> false, " +
				"computed: false -> true, " +
				"default: <nil> -> 10, " +
				"sensitive: false -> true, " +
				`conflicts with: [] -> [a], ` +
				`required with: [a] -> [b], ` +
				`exactly one of: [] -> [a], ` +
				`at least one of: [] -> []`,
		},
		{
			name: "Provider config block add",
			change: BlockChange{
				Scope: ProviderScope{},
				Path:  []string{"foo", "bar"},
				IsAdd: true,
			},
			expect: `Block "foo.bar" of provider config is added`,
		},
		{
			name: "Resource block add",
			change: BlockChange{
				Scope: ResourceScope{Type: "foo_resource"},
				Path:  []string{"foo", "bar"},
				IsAdd: true,
			},
			expect: `Block "foo.bar" of resource foo_resource is added`,
		},
		{
			name: "Data Source block add",
			change: BlockChange{
				Scope: DataSourceScope{Type: "foo_resource"},
				Path:  []string{"foo", "bar"},
				IsAdd: true,
			},
			expect: `Block "foo.bar" of data source foo_resource is added`,
		},
		{
			name: "Resource block delete",
			change: BlockChange{
				Scope:    ResourceScope{Type: "foo_resource"},
				Path:     []string{"foo", "bar"},
				IsDelete: true,
			},
			expect: `Block "foo.bar" of resource foo_resource is deleted`,
		},
		{
			name: "Block attribute change (complete)",
			change: BlockChange{
				Scope:    ResourceScope{Type: "foo_resource"},
				Path:     []string{"foo", "bar"},
				IsModify: true,
				Modification: &BlockModify{
					NestingMode: &Modification[schema.NestingMode]{
						From: schema.NestingSingle,
						To:   schema.NestingGroup,
					},
					Required: &Modification[bool]{
						From: false,
						To:   true,
					},
					Optional: &Modification[bool]{
						From: true,
						To:   false,
					},
					Computed: &Modification[bool]{
						From: false,
						To:   true,
					},
					// Intentionally ignore this
					// ForceNew: &Modification[bool]{
					// 	From: false,
					// 	To:   false,
					// },
					ConflictsWith: &Modification[[]string]{
						From: []string{},
						To:   []string{"a"},
					},
					RequiredWith: &Modification[[]string]{
						From: []string{"a"},
						To:   []string{"b"},
					},
					ExactlyOneOf: &Modification[[]string]{
						From: nil,
						To:   []string{"a"},
					},
					AtLeastOneOf: &Modification[[]string]{
						From: []string{},
						To:   nil,
					},
					MinItems: &Modification[int]{
						From: 0,
						To:   1,
					},
					MaxItems: &Modification[int]{
						From: 1,
						To:   2,
					},
				},
			},
			expect: `Block "foo.bar" of resource foo_resource is changed: ` +
				"nesting mode: 1 -> 2, " +
				"required: false -> true, " +
				"optional: true -> false, " +
				"computed: false -> true, " +
				`conflicts with: [] -> [a], ` +
				`required with: [a] -> [b], ` +
				`exactly one of: [] -> [a], ` +
				`at least one of: [] -> [], ` +
				`min items: 0 -> 1, ` +
				`max items: 1 -> 2`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.change.String())
		})
	}
}
