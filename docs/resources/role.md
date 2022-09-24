---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "opensearch_role Resource - terraform-provider-opensearch"
subcategory: ""
description: |-
  Opensearch role.
---

# opensearch_role (Resource)

Opensearch role.

## Example Usage

```terraform
resource "opensearch_role" "qa" {
  name = "qa"

  index_permissions {
      index_patterns = ["qa*"]
      allowed_actions = ["read", "write"]
  }

  cluster_permissions = [
    "indices_monitor"
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) Name of the role.

### Optional

- **cluster_permissions** (Set of String) Permissions for cluster wide actions the role has.
- **id** (String) The ID of this resource.
- **index_permissions** (Block Set) Permissions for index access the role has. (see [below for nested schema](#nestedblock--index_permissions))
- **tenant_permissions** (Block Set) Permissions for tenant access the role has. (see [below for nested schema](#nestedblock--tenant_permissions))

<a id="nestedblock--index_permissions"></a>
### Nested Schema for `index_permissions`

Optional:

- **allowed_actions** (Set of String)
- **document_level_security** (String)
- **field_level_security** (Set of String)
- **index_patterns** (Set of String)
- **masked_fields** (Set of String)


<a id="nestedblock--tenant_permissions"></a>
### Nested Schema for `tenant_permissions`

Optional:

- **allowed_actions** (Set of String)
- **tenant_patterns** (Set of String)

