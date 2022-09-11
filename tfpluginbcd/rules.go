package tfpluginbcd

type Rule struct {
	ID          string
	Description string
	Expr        string
}

var Rules = map[string]Rule{
	"R001": {
		ID:          "R001",
		Description: "A resource is deleted",
		Expr:        `c.kind == "resource"; not c.is_data_source; c.is_delete`,
	},
	"R002": {
		ID:          "R002",
		Description: "A data source is deleted",
		Expr:        `c.kind == "resource"; c.is_data_source; c.is_delete`,
	},
	"R003": {
		ID:          "R003",
		Description: "An attribute is deleted",
		Expr:        `c.kind == "attribute"; c.is_delete`,
	},
	"R004": {
		ID:          "R004",
		Description: "A block is deleted",
		Expr:        `c.kind == "block"; c.is_delete`,
	},
	"R005": {
		ID:          "R005",
		Description: "The type of an attribute is changed",
		Expr:        `c.kind == "attribute"; c.is_modify; c.modification.type`,
	},
	"R006": {
		ID:          "R006",
		Description: "An optional attribute is changed to be required",
		Expr:        `c.kind == "attribute"; c.is_modify; c.modification.required.to == true`,
	},
	"R007": {
		ID:          "R007",
		Description: "An optional block is changed to be required",
		Expr:        `c.kind == "block"; c.is_modify; c.modification.required.to == true`,
	},
	"R008": {
		ID:          "R008",
		Description: "A new required attribute is added",
		Expr:        `c.kind == "attribute"; c.is_add; c.current.required == true`,
	},
	"R009": {
		ID:          "R009",
		Description: "A new required block is added",
		Expr:        `c.kind == "block"; c.is_add; c.current.required == true`,
	},
}
