---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "netbox_tenancy_tenant Resource - terraform-provider-netbox"
subcategory: ""
description: |-
  Manage a tenant (tenancy module) within Netbox.
---

# netbox_tenancy_tenant (Resource)

Manage a tenant (tenancy module) within Netbox.

## Example Usage

```terraform
resource "netbox_tenancy_tenant" "tenant_test" {
  name            = "TestTenant"
  slug            = "TestTenant"
  description     = "Tenant created by terraform"
  comments        = "Some test comments"
  tenant_group_id = netbox_tenancy_tenant_group.tenant_group_test.id

  tag {
    name = "tag1"
    slug = "tag1"
  }

  custom_field {
    name = "cf_boolean"
    type = "boolean"
    value = "true"
  }

  custom_field {
    name = "cf_date"
    type = "date"
    value = "2020-12-25"
  }

  custom_field {
    name = "cf_text"
    type = "text"
    value = "some text"
  }

  custom_field {
    name = "cf_integer"
    type = "integer"
    value = "10"
  }

  custom_field {
    name = "cf_selection"
    type = "select"
    value = "1"
  }

  custom_field {
    name = "cf_url"
    type = "url"
    value = "https://github.com"
  }

  custom_field {
    name = "cf_multi_selection"
    type = "multiselect"
    value = jsonencode([
      "0",
      "1"
    ])
  }

  custom_field {
    name = "cf_json"
    type = "json"
    value = jsonencode({
      stringvalue = "string"
      boolvalue = false
      dictionary = {
        numbervalue = 5
      }
    })
  }

  custom_field {
    name = "cf_object"
    type = "object"
    value = data.netbox_dcim_platform.platform_test.id
  }

  custom_field {
    name = "cf_multi_object"
    type = "multiobject"
    value = jsonencode([
      data.netbox_dcim_platform.platform_test.id,
      data.netbox_dcim_platform.platform_test2.id
    ])
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name for this tenant (tenancy module).
- `slug` (String) The slug for this tenant (tenancy module).

### Optional

- `comments` (String) Comments for this tenant (tenancy module).
- `custom_field` (Block Set) Existing custom fields to associate to this ressource. (see [below for nested schema](#nestedblock--custom_field))
- `description` (String) The description for this tenant (tenancy module).
- `tag` (Block Set) Existing tag to associate to this resource. (see [below for nested schema](#nestedblock--tag))
- `tenant_group_id` (Number) ID of the group where this tenant (tenancy module) is attached to.

### Read-Only

- `content_type` (String) The content type of this tenant (tenancy module).
- `id` (String) The ID of this resource.

<a id="nestedblock--custom_field"></a>
### Nested Schema for `custom_field`

Required:

- `name` (String) Name of the existing custom field.
- `type` (String) Type of the existing custom field (text, longtext, integer, boolean, date, url, json, select, multiselect, object, multiobject, selection (deprecated), multiple(deprecated)).
- `value` (String) Value of the existing custom field.


<a id="nestedblock--tag"></a>
### Nested Schema for `tag`

Required:

- `name` (String) Name of the existing tag.
- `slug` (String) Slug of the existing tag.


