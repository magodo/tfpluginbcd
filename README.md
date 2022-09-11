# tfpluginbcd

Terraform plugin breaking change detector.

## Install

`go install github.com/magodo/tfpluginbcd@latest`

## Usage

`tfpluginbcd` detects breaking changes between two [terraform plugin schemas](https://github.com/magodo/tfpluginschema). The typical workflow is as below (given your terraform plugin is based on [terraform-plugin-sdk](https://github.com/hashicorp/terraform-plugin-sdk)):

- Checkout the terraform plugin project with version `v1`, run the helper script from the project root dir: `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/magodo/tfpluginschema/main/tool/schema_dumper_sdk_v2/run.sh)" bash <your_provider_func>` (Replace `your_provider_func` with your provider's init function in form of `<package path>.<function name>`. E.g. for [terraform-provider-azurerm](https://github.com/hashicorp/terraform-provider-azurerm), it is `github.com/hashicorp/terraform-provider-azurerm/internal/provider.AzureProvider`). Redirect the output schema to a file called *schema_v1.json*
- Repeat above for `v2`, output the schema to a file called *schema_v2.json*
- Run `tfpluginbcd run -all schema_v1.json schema_v2.json` (You can also select a subset of rules by `--rule` option, or feed your custom rules via `--custom-rule`) to show any breaking change between `v1` and `v2`

## Rules

### Pre-defined Rules

`tfpluginbcd` defines several rules which are regarded as breaking changes for most of users:

|Name|Description|Rego Expression|
|-|-|-|
|R001|A resource is deleted|c.kind == "resource"; not c.is_data_source; c.is_delete|
|R002|A data source is deleted|c.kind == "resource"; c.is_data_source; c.is_delete|
|R003|An attribute is deleted|c.kind == "attribute"; c.is_delete|
|R004|A block is deleted|c.kind == "block"; c.is_delete|
|R005|The type of an attribute is changed|c.kind == "attribute"; c.is_modify; c.modification.type|
|R006|An optional attribute is changed to be required|c.kind == "attribute"; c.is_modify; c.modification.required.to == true|
|R007|An optional block is changed to be required|c.kind == "block"; c.is_modify; c.modification.required.to == true|
|R008|A new required attribute is added|c.kind == "attribute"; c.is_add; c.current.required == true|
|R009|A new required block is added|c.kind == "block"; c.is_add; c.current.required == true|

### Custom Rules

Users can specify custom rules via the `--custom-rule` option. `tfpluginbcd` uses [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/) to define the breaking change rules. The content of the `--custom-rule` is [Rego expressions](https://www.openpolicyagent.org/docs/latest/policy-language/#multiple-expressions), where users are provided with a special reference `c` that represents each schema change.

The definition of the schema change (i.e. `c`) can be one of below:

1. Resource/DataSource Change: Resource/Data Source level schema changes

    ```
    {
        "kind"          : "resource",
        "type"          : string,                   # The terraform resource type
        "is_data_source": bool,

        # Exactly one of below can be true
        "is_add"        : bool,
        "is_delete"     : bool,
        "is_modify"     : bool,

        "current"       : <Resource>,               # The current resource schema, which is present only when is_add/is_modify is true
        "modification"  : <ResourceModification>    # The resource schema modification, which present only when is_modify is true
    }
    ```

    The `Resource` is defined as:

    ```
    {
        "schema_version": int # The resource's schema version
    }
    ```

    The `ResourceModification` has the same fields as `Resource`, except each field is a `Modification` object, which is present only when that field is changed.

2. Attribute Change: Attribute level schema changes, which includes provider, resource and data source attributes

    ```
    {
        "kind"          : "attribute",
        "scope"         : <Scope>,                  # The scope of this attribute
        "path"          : []string                  # The path to the attribute

        # Exactly one of below can be true
        "is_add"        : bool,
        "is_delete"     : bool,
        "is_modify"     : bool,

        "current"       : <Resource>,               # The current attribute schema, which is present only when is_add/is_modify is true
        "modification"  : <ResourceModification>    # The attribute schema modification, which present only when is_modify is true
    }
    ```

    The `Attribute`  is defined as:

    ```
    {
        "kind"              : "attribute",
        "type"              : cty.Type,
        "required"          : bool,
        "optional"          : bool,
        "computed"          : bool,
        "force_new"         : bool,
        "default"           : any,
        "sensitive"         : bool,
        "conflicts_with"    : []string,
        "required_with"     : []string,
        "at_leatst_one_of"  : []string,
        "exactly_one_of"    : []string
    }
    ```

    The `AttributeModification` has the same fields as `Attribute`, except each field is a `Modification` object, which is present only when that field is changed.

3. Block Change: Block level schema changes, which includes provider, resource and data source blocks

    ```
    {
        "kind"          : "block",
        "scope"         : <Scope>,                  # The scope of this attribute
        "path"          : []string                  # The path to the attribute

        # Exactly one of below can be true
        "is_add"        : bool,
        "is_delete"     : bool,
        "is_modify"     : bool,

        "current"       : <Block>,                  # The current block schema, which is present only when is_add/is_modify is true
        "modification"  : <BlockModification>       # The block schema modification, which present only when is_modify is true
    }
    ```

    The `Block`  is defined as:

    ```
    {
        "kind"              : "block",
        "nesting_mode"      : int,
        "required"          : bool,
        "optional"          : bool,
        "computed"          : bool,
        "force_new"         : bool,
        "conflicts_with"    : []string,
        "required_with"     : []string,
        "at_leatst_one_of"  : []string,
        "exactly_one_of"    : []string,
        "min_items"         : int,
        "max_items"         : int
    }
    ```

    The `BlockModification` has the same fields as `Block`, except each field is a `Modification` object, which is present only when that field is changed.

Additionally:

- The `Modification` object is defined as:

    ```
    {
        From: any,  # The source value
        To: any     # The destination value
    }
    ```

- The `Scope` object can be one of below:

    - Provider scope:

        ```
        {
            "kind": "provider"
        }
        ```
    - Reosurce/DataSource scope:

        ```
        {
            "kind"          : "resource",
            "type"          : string,       # The terraform resource type
            "is_data_source": bool
        }
        ```

Examples custom rules:

|Description|Rego Expression|
|-|-|
|Set a default value to an attribute (which was `null`)| `c.kind == "attribute"; c.modification["default"].from == null`|
|The `max_items` is decreased for a block| `c.kind == "block"; c.modification.max_items.to < c.modification.max_items.from`|