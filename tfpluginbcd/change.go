package tfpluginbcd

import (
	"fmt"
	"strings"

	"github.com/magodo/tfpluginschema/schema"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/exp/slices"
)

type Change interface {
	isChange()
	String() string
}

type ProviderChange struct {
	// Exactly one of them is true
	IsAdd    bool
	IsDelete bool
}

func (ProviderChange) isChange() {}

func (c ProviderChange) String() string {
	var verb string
	switch {
	case c.IsAdd:
		verb = "added"
	case c.IsDelete:
		verb = "deleted"
	}
	return fmt.Sprintf("Provider config is %s", verb)
}

type ResourceChange struct {
	// Resource type
	Type         string
	IsDataSource bool

	// Exactly one of them is true
	IsAdd    bool
	IsDelete bool
	IsModify bool

	// Current represents the current schema of this resource, it is nil if IsDelete is true.
	Current *Resource

	// Modification represents the modification of this resource, it is non-nil only when IsModify is true.
	Modification *ResourceModify
}

func (ResourceChange) isChange() {}

func (c ResourceChange) String() string {
	var msg string

	if c.IsDataSource {
		msg += "Data Source"
	} else {
		msg += "Resource"
	}

	msg += " " + c.Type + " is"

	switch {
	case c.IsAdd:
		msg += " added"
	case c.IsDelete:
		msg += " deleted"
	case c.IsModify:
		msg += " changed: " + c.Modification.String()
	}
	return msg
}

type Scope interface {
	isScope()
}

type ProviderScope struct{}

func (ProviderScope) isScope() {}

type ResourceScope struct {
	Type         string
	IsDataSource bool
}

func (ResourceScope) isScope() {}

type AttributeChange struct {
	Scope
	Path []string

	// Exactly one of them is true
	IsAdd    bool
	IsDelete bool
	IsModify bool

	// Current represents the current schema of this attribute, it is nil if IsDelete is true.
	Current *Attribute

	// Modification represents the modification of this attribute, it is non-nil only when IsModify is true.
	Modification *AttributeModify
}

func (AttributeChange) isChange() {}

func (c AttributeChange) String() string {
	var msg string

	msg += fmt.Sprintf("Attribute %q of", strings.Join(c.Path, "."))

	switch scope := c.Scope.(type) {
	case ProviderScope:
		msg += " provider config"
	case ResourceScope:
		if scope.IsDataSource {
			msg += " data source " + scope.Type
		} else {
			msg += " resource " + scope.Type
		}
	}

	msg += " is"

	switch {
	case c.IsAdd:
		msg += " added"
	case c.IsDelete:
		msg += " deleted"
	case c.IsModify:
		msg += " changed: " + c.Modification.String()
	}

	return msg
}

type BlockChange struct {
	Scope
	Path []string

	// Exactly one of them is true
	IsAdd    bool
	IsDelete bool
	IsModify bool

	// Current represents the current schema of this block, it is nil if IsDelete is true.
	Current *Block

	// Modification represents the modification of this block, it is non-nil only when IsModify is true.
	Modification *BlockModify
}

func (BlockChange) isChange() {}

func (c BlockChange) String() string {
	var msg string

	msg += fmt.Sprintf("Block %q of", strings.Join(c.Path, "."))

	switch scope := c.Scope.(type) {
	case ProviderScope:
		msg += " provider config"
	case ResourceScope:
		if scope.IsDataSource {
			msg += " data source " + scope.Type
		} else {
			msg += " resource " + scope.Type
		}
	}

	msg += " is"

	switch {
	case c.IsAdd:
		msg += " added"
	case c.IsDelete:
		msg += " deleted"
	case c.IsModify:
		msg += " changed: " + c.Modification.String()
	}

	return msg
}

type Modification[T any] struct {
	From T
	To   T
}

type Resource struct {
	SchemaVersion int
}

type ResourceModify struct {
	SchemaVersion *Modification[int]
}

func (m ResourceModify) String() string {
	var l []string
	if m.SchemaVersion != nil {
		l = append(l, fmt.Sprintf("schema version: %d -> %d", m.SchemaVersion.From, m.SchemaVersion.To))
	}
	return strings.Join(l, ", ")
}

type Attribute struct {
	Type          cty.Type
	Required      bool
	Optional      bool
	Computed      bool
	ForceNew      bool
	Default       interface{}
	Sensitive     bool
	ConflictsWith []string
	ExactlyOneOf  []string
	AtLeastOneOf  []string
	RequiredWith  []string
}

type AttributeModify struct {
	Type          *Modification[cty.Type]
	Required      *Modification[bool]
	Optional      *Modification[bool]
	Computed      *Modification[bool]
	ForceNew      *Modification[bool]
	Default       *Modification[any]
	Sensitive     *Modification[bool]
	ConflictsWith *Modification[[]string]
	RequiredWith  *Modification[[]string]
	ExactlyOneOf  *Modification[[]string]
	AtLeastOneOf  *Modification[[]string]
}

func (m AttributeModify) String() string {
	var l []string
	if m.Type != nil {
		l = append(l, fmt.Sprintf("type: %s -> %s", m.Type.From.FriendlyName(), m.Type.To.FriendlyName()))
	}
	if m.Required != nil {
		l = append(l, fmt.Sprintf("required: %t -> %t", m.Required.From, m.Required.To))
	}
	if m.Optional != nil {
		l = append(l, fmt.Sprintf("optional: %t -> %t", m.Optional.From, m.Optional.To))
	}
	if m.Computed != nil {
		l = append(l, fmt.Sprintf("computed: %t -> %t", m.Computed.From, m.Computed.To))
	}
	if m.ForceNew != nil {
		l = append(l, fmt.Sprintf("force new: %t -> %t", m.ForceNew.From, m.ForceNew.To))
	}
	if m.Default != nil {
		l = append(l, fmt.Sprintf("default: %v -> %v", m.Default.From, m.Default.To))
	}
	if m.Sensitive != nil {
		l = append(l, fmt.Sprintf("sensitive: %t -> %t", m.Sensitive.From, m.Sensitive.To))
	}
	if m.ConflictsWith != nil {
		l = append(l, fmt.Sprintf("conflicts with: [%s] -> [%s]", strings.Join(m.ConflictsWith.From, ", "), strings.Join(m.ConflictsWith.To, ", ")))
	}
	if m.RequiredWith != nil {
		l = append(l, fmt.Sprintf("required with: [%s] -> [%s]", strings.Join(m.RequiredWith.From, ", "), strings.Join(m.RequiredWith.To, ", ")))
	}
	if m.ExactlyOneOf != nil {
		l = append(l, fmt.Sprintf("exactly one of: [%s] -> [%s]", strings.Join(m.ExactlyOneOf.From, ", "), strings.Join(m.ExactlyOneOf.To, ", ")))
	}
	if m.AtLeastOneOf != nil {
		l = append(l, fmt.Sprintf("at least one of: [%s] -> [%s]", strings.Join(m.AtLeastOneOf.From, ", "), strings.Join(m.AtLeastOneOf.To, ", ")))
	}
	return strings.Join(l, ", ")
}

type Block struct {
	NestingMode   schema.NestingMode
	Required      bool
	Optional      bool
	Computed      bool
	ForceNew      bool
	ConflictsWith []string
	ExactlyOneOf  []string
	AtLeastOneOf  []string
	RequiredWith  []string
	MinItems      int
	MaxItems      int
}

type BlockModify struct {
	NestingMode   *Modification[schema.NestingMode]
	Required      *Modification[bool]
	Optional      *Modification[bool]
	Computed      *Modification[bool]
	ForceNew      *Modification[bool]
	ConflictsWith *Modification[[]string]
	ExactlyOneOf  *Modification[[]string]
	AtLeastOneOf  *Modification[[]string]
	RequiredWith  *Modification[[]string]
	MinItems      *Modification[int]
	MaxItems      *Modification[int]
}

func (m BlockModify) String() string {
	var l []string
	if m.NestingMode != nil {
		l = append(l, fmt.Sprintf("nesting mode: %v -> %v", m.NestingMode.From, m.NestingMode.To))
	}
	if m.Required != nil {
		l = append(l, fmt.Sprintf("required: %t -> %t", m.Required.From, m.Required.To))
	}
	if m.Optional != nil {
		l = append(l, fmt.Sprintf("optional: %t -> %t", m.Optional.From, m.Optional.To))
	}
	if m.Computed != nil {
		l = append(l, fmt.Sprintf("computed: %t -> %t", m.Computed.From, m.Computed.To))
	}
	if m.ForceNew != nil {
		l = append(l, fmt.Sprintf("force new: %t -> %t", m.ForceNew.From, m.ForceNew.To))
	}
	if m.ConflictsWith != nil {
		l = append(l, fmt.Sprintf("conflicts with: [%s] -> [%s]", strings.Join(m.ConflictsWith.From, ", "), strings.Join(m.ConflictsWith.To, ", ")))
	}
	if m.RequiredWith != nil {
		l = append(l, fmt.Sprintf("required with: [%s] -> [%s]", strings.Join(m.RequiredWith.From, ", "), strings.Join(m.RequiredWith.To, ", ")))
	}
	if m.ExactlyOneOf != nil {
		l = append(l, fmt.Sprintf("exactly one of: [%s] -> [%s]", strings.Join(m.ExactlyOneOf.From, ", "), strings.Join(m.ExactlyOneOf.To, ", ")))
	}
	if m.AtLeastOneOf != nil {
		l = append(l, fmt.Sprintf("at least one of: [%s] -> [%s]", strings.Join(m.AtLeastOneOf.From, ", "), strings.Join(m.AtLeastOneOf.To, ", ")))
	}
	if m.MinItems != nil {
		l = append(l, fmt.Sprintf("min items: %d -> %d", m.MinItems.From, m.MinItems.To))
	}
	if m.MaxItems != nil {
		l = append(l, fmt.Sprintf("max items: %d -> %d", m.MaxItems.From, m.MaxItems.To))
	}
	return strings.Join(l, ", ")
}

func NewAttribute(attr *schema.Attribute) *Attribute {
	if attr == nil {
		return nil
	}
	return &Attribute{
		Type:          attr.Type,
		Required:      attr.Required,
		Optional:      attr.Optional,
		Computed:      attr.Computed,
		ForceNew:      attr.ForceNew,
		Default:       attr.Default,
		Sensitive:     attr.Sensitive,
		ConflictsWith: attr.ConflictsWith,
		ExactlyOneOf:  attr.ExactlyOneOf,
		AtLeastOneOf:  attr.AtLeastOneOf,
		RequiredWith:  attr.RequiredWith,
	}
}

func NewNestedBlock(blk *schema.NestedBlock) *Block {
	if blk == nil {
		return nil
	}
	return &Block{
		NestingMode:   blk.NestingMode,
		Required:      blk.Required,
		Optional:      blk.Optional,
		Computed:      blk.Computed,
		ForceNew:      blk.ForceNew,
		ConflictsWith: blk.ConflictsWith,
		ExactlyOneOf:  blk.ExactlyOneOf,
		AtLeastOneOf:  blk.AtLeastOneOf,
		RequiredWith:  blk.RequiredWith,
		MinItems:      blk.MinItems,
		MaxItems:      blk.MaxItems,
	}
}

func NewAttributeModify(oattr schema.Attribute, nattr schema.Attribute) *AttributeModify {
	isChanged := false
	ret := &AttributeModify{}

	if !oattr.Type.Equals(nattr.Type) {
		isChanged = true
		ret.Type = &Modification[cty.Type]{
			From: oattr.Type,
			To:   nattr.Type,
		}
	}
	if oattr.Required != nattr.Required {
		isChanged = true
		ret.Required = &Modification[bool]{
			From: oattr.Required,
			To:   nattr.Required,
		}
	}
	if oattr.Optional != nattr.Optional {
		isChanged = true
		ret.Optional = &Modification[bool]{
			From: oattr.Optional,
			To:   nattr.Optional,
		}
	}
	if oattr.Computed != nattr.Computed {
		isChanged = true
		ret.Computed = &Modification[bool]{
			From: oattr.Computed,
			To:   nattr.Computed,
		}
	}
	if oattr.ForceNew != nattr.ForceNew {
		isChanged = true
		ret.ForceNew = &Modification[bool]{
			From: oattr.ForceNew,
			To:   nattr.ForceNew,
		}
	}
	if oattr.Default != nattr.Default {
		isChanged = true
		ret.Default = &Modification[any]{
			From: oattr.Default,
			To:   nattr.Default,
		}
	}
	if oattr.Sensitive != nattr.Sensitive {
		isChanged = true
		ret.Sensitive = &Modification[bool]{
			From: oattr.Sensitive,
			To:   nattr.Sensitive,
		}
	}
	if !slices.Equal(oattr.ConflictsWith, nattr.ConflictsWith) {
		isChanged = true
		ret.ConflictsWith = &Modification[[]string]{
			From: oattr.ConflictsWith,
			To:   nattr.ConflictsWith,
		}
	}
	if !slices.Equal(oattr.RequiredWith, nattr.RequiredWith) {
		isChanged = true
		ret.RequiredWith = &Modification[[]string]{
			From: oattr.RequiredWith,
			To:   nattr.RequiredWith,
		}
	}
	if !slices.Equal(oattr.AtLeastOneOf, nattr.AtLeastOneOf) {
		isChanged = true
		ret.AtLeastOneOf = &Modification[[]string]{
			From: oattr.AtLeastOneOf,
			To:   nattr.AtLeastOneOf,
		}
	}
	if !slices.Equal(oattr.ExactlyOneOf, nattr.ExactlyOneOf) {
		isChanged = true
		ret.ExactlyOneOf = &Modification[[]string]{
			From: oattr.ExactlyOneOf,
			To:   nattr.ExactlyOneOf,
		}
	}
	if !isChanged {
		return nil
	}
	return ret
}

func NewBlockModify(oblk schema.NestedBlock, nblk schema.NestedBlock) *BlockModify {
	isChanged := false
	ret := &BlockModify{}

	if oblk.NestingMode != nblk.NestingMode {
		isChanged = true
		ret.NestingMode = &Modification[schema.NestingMode]{
			From: oblk.NestingMode,
			To:   nblk.NestingMode,
		}
	}
	if oblk.Required != nblk.Required {
		isChanged = true
		ret.Required = &Modification[bool]{
			From: oblk.Required,
			To:   nblk.Required,
		}
	}
	if oblk.Optional != nblk.Optional {
		isChanged = true
		ret.Optional = &Modification[bool]{
			From: oblk.Optional,
			To:   nblk.Optional,
		}
	}
	if oblk.Computed != nblk.Computed {
		isChanged = true
		ret.Computed = &Modification[bool]{
			From: oblk.Computed,
			To:   nblk.Computed,
		}
	}
	if oblk.ForceNew != nblk.ForceNew {
		isChanged = true
		ret.ForceNew = &Modification[bool]{
			From: oblk.ForceNew,
			To:   nblk.ForceNew,
		}
	}
	if !slices.Equal(oblk.ConflictsWith, nblk.ConflictsWith) {
		isChanged = true
		ret.ConflictsWith = &Modification[[]string]{
			From: oblk.ConflictsWith,
			To:   nblk.ConflictsWith,
		}
	}
	if !slices.Equal(oblk.RequiredWith, nblk.RequiredWith) {
		isChanged = true
		ret.RequiredWith = &Modification[[]string]{
			From: oblk.RequiredWith,
			To:   nblk.RequiredWith,
		}
	}
	if !slices.Equal(oblk.AtLeastOneOf, nblk.AtLeastOneOf) {
		isChanged = true
		ret.AtLeastOneOf = &Modification[[]string]{
			From: oblk.AtLeastOneOf,
			To:   nblk.AtLeastOneOf,
		}
	}
	if !slices.Equal(oblk.ExactlyOneOf, nblk.ExactlyOneOf) {
		isChanged = true
		ret.ExactlyOneOf = &Modification[[]string]{
			From: oblk.ExactlyOneOf,
			To:   nblk.ExactlyOneOf,
		}
	}
	if oblk.MinItems != nblk.MinItems {
		isChanged = true
		ret.MinItems = &Modification[int]{
			From: oblk.MinItems,
			To:   nblk.MinItems,
		}
	}
	if oblk.MaxItems != nblk.MaxItems {
		isChanged = true
		ret.MaxItems = &Modification[int]{
			From: oblk.MaxItems,
			To:   nblk.MaxItems,
		}
	}
	if !isChanged {
		return nil
	}
	return ret
}
