---
page_title: "auxo_contact Data Source - terraform-provider-auxo"
subcategory: ""
description: |-
  A contact which can be used a.o. as main- or securitycontact in a protectsurface.
---

# auxo_contact (Data Source)

A contact which can be used a.o. as main- or securitycontact in a `protectsurface`.

## Example Usage

```terraform
data "auxo_contact" "rob" {
  email = "rob.maas+tst@on2it.net"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `email` (String) Emails of the contact

### Read-Only

- `id` (String) Computed unique IDs of the contact
