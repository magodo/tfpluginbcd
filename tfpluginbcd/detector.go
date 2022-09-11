package tfpluginbcd

import (
	"github.com/magodo/tfpluginschema/schema"
	"golang.org/x/exp/slices"
)

func Compare(oldSch, newSch *schema.ProviderSchema) []Change {
	var changes []Change

	switch {
	case (oldSch.Provider != nil && oldSch.Provider.Block != nil) && (newSch.Provider != nil && newSch.Provider.Block != nil):
		changes = append(changes, compareBlock(ProviderScope{}, []string{}, oldSch.Provider.Block, newSch.Provider.Block)...)
	case (oldSch.Provider == nil || oldSch.Provider.Block == nil) && (newSch.Provider == nil || newSch.Provider.Block == nil):
		// do nothing
	case (oldSch.Provider == nil || oldSch.Provider.Block == nil) && (newSch.Provider != nil && newSch.Provider.Block != nil):
		changes = append(changes, ProviderChange{IsAdd: true})
	case (oldSch.Provider != nil && oldSch.Provider.Block != nil) && (newSch.Provider == nil || newSch.Provider.Block == nil):
		changes = append(changes, ProviderChange{IsDelete: true})
	}

	changes = append(changes, compareResources(oldSch.DataSourceSchemas, newSch.DataSourceSchemas, true)...)
	changes = append(changes, compareResources(oldSch.ResourceSchemas, newSch.ResourceSchemas, false)...)

	return changes
}

func compareResources(orm, nrm map[string]*schema.Resource, isDataSource bool) []Change {
	var changes []Change
	for _, rt := range mapSortedKeys(orm) {
		ores := orm[rt]
		nres, ok := nrm[rt]
		// Delete
		if !ok {
			changes = append(changes, ResourceChange{
				Type:         rt,
				IsDataSource: isDataSource,
				IsDelete:     true,
			})
			continue
		}
		// Update
		if ores.SchemaVersion != nres.SchemaVersion {
			changes = append(changes, ResourceChange{
				Type:         rt,
				IsDataSource: isDataSource,
				IsModify:     true,
				Current: &Resource{
					SchemaVersion: nres.SchemaVersion,
				},
				Modification: &ResourceModify{
					SchemaVersion: &Modification[int]{
						From: ores.SchemaVersion,
						To:   nres.SchemaVersion,
					},
				},
			})
		}
		// Inner
		scope := ResourceScope{
			Type:         rt,
			IsDataSource: isDataSource,
		}
		changes = append(changes, compareBlock(scope, []string{}, ores.Block, nres.Block)...)
	}

	for _, rt := range mapSortedKeys(nrm) {
		nres := nrm[rt]
		// Add
		if _, ok := orm[rt]; !ok {
			changes = append(changes, ResourceChange{
				Type:         rt,
				IsDataSource: isDataSource,
				IsAdd:        true,
				Current: &Resource{
					SchemaVersion: nres.SchemaVersion,
				},
			})
			continue
		}
	}

	return changes
}

func compareBlock(scope Scope, path []string, oblk, nblk *schema.Block) []Change {
	var changes []Change

	for _, name := range mapSortedKeys(oblk.Attributes) {
		oattr := oblk.Attributes[name]
		changes = append(changes, compareAttribute(scope, append(path, name), oattr, nblk.Attributes[name])...)
	}
	for _, name := range mapSortedKeys(nblk.Attributes) {
		nattr := nblk.Attributes[name]
		if _, ok := oblk.Attributes[name]; ok {
			continue
		}
		changes = append(changes, compareAttribute(scope, append(path, name), nil, nattr)...)
	}

	for _, name := range mapSortedKeys(oblk.NestedBlocks) {
		oNestBlk := oblk.NestedBlocks[name]
		changes = append(changes, compareNestedBlock(scope, append(path, name), oNestBlk, nblk.NestedBlocks[name])...)
	}
	for _, name := range mapSortedKeys(nblk.NestedBlocks) {
		nNestBlk := nblk.NestedBlocks[name]
		if _, ok := oblk.NestedBlocks[name]; ok {
			continue
		}
		changes = append(changes, compareNestedBlock(scope, append(path, name), nil, nNestBlk)...)
	}

	return changes
}

func compareAttribute(scope Scope, path []string, oattr, nattr *schema.Attribute) []Change {
	if nattr == nil {
		return []Change{
			AttributeChange{
				Scope:    scope,
				Path:     path,
				IsDelete: true,
			},
		}
	}
	if oattr == nil {
		return []Change{
			AttributeChange{
				Scope:   scope,
				Path:    path,
				IsAdd:   true,
				Current: NewAttribute(nattr),
			},
		}
	}

	modification := NewAttributeModify(*oattr, *nattr)
	if modification == nil {
		return nil
	}
	return []Change{
		AttributeChange{
			Scope:        scope,
			Path:         path,
			IsModify:     true,
			Current:      NewAttribute(nattr),
			Modification: modification,
		},
	}
}

func compareNestedBlock(scope Scope, path []string, oblk, nblk *schema.NestedBlock) []Change {
	if nblk == nil {
		return []Change{
			BlockChange{
				Scope:    scope,
				Path:     path,
				IsDelete: true,
			},
		}
	}
	if oblk == nil {
		return []Change{
			BlockChange{
				Scope:   scope,
				Path:    path,
				IsAdd:   true,
				Current: NewNestedBlock(nblk),
			},
		}
	}

	var changes []Change
	modification := NewBlockModify(*oblk, *nblk)
	if modification != nil {
		changes = append(changes, BlockChange{
			Scope:        scope,
			Path:         path,
			IsModify:     true,
			Current:      NewNestedBlock(nblk),
			Modification: modification,
		})
	}

	changes = append(changes, compareBlock(scope, path, oblk.Block, nblk.Block)...)
	return changes
}

func mapSortedKeys[T any](m map[string]T) []string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}
