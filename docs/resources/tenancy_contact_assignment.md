---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "netbox_tenancy_contact_assignment Resource - terraform-provider-netbox"
subcategory: ""
description: |-
  Link a contact (tenancy module) to another resource within Netbox.
---

# netbox_tenancy_contact_assignment (Resource)

Link a contact (tenancy module) to another resource within Netbox.

## Example Usage

```terraform
resource "netbox_tenancy_contact_assignment" "contact_assignment_01" {
  contact_id = netbox_tenancy_contact.contact.id
  contact_role_id = netbox_tenancy_contact_role.contact_role_02.id
  content_type = netbox_virtualization_vm.vm_test.content_type
  object_id = netbox_virtualization_vm.vm_test.id
  priority = "primary"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `contact_id` (Number) ID of the contact to link to this contact assignment (tenancy module).
- `contact_role_id` (Number) The role of the contact for this contact assignment (tenancy module).
- `content_type` (String) Type of the object where the contact will be linked.
- `object_id` (Number) ID of the object where the contact will be linked.

### Optional

- `priority` (String) Priority of this contact among primary, secondary and tertiary (primary by default).

### Read-Only

- `id` (String) The ID of this resource.

