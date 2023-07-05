---
page_title: "auxo_state Resource - terraform-provider-auxo"
subcategory: ""
description: |-
  A state contains resources and their location, belonging to a protect surface.
---

# auxo_state (Resource)

A state contains resources and their location, belonging to a protect surface.

## Example Usage

```terraform
resource "auxo_state" "ps_ad-loc_zaltbommel-ipv4" {
  content_type      = "ipv4"
  description       = "IPv4 allocations of AD servers"
  location_id       = auxo_location.loc_zaltbommel.id
  protectsurface_id = auxo_protectsurface.ps_ad.id
  content           = ["10.0.42.10", "10.0.42.11", "10.0.42.12"]
}
```

Current supported `content_type` are:

| content type  | description                                                                                      |
| ------------- | ------------------------------------------------------------------------------------------------ |
| azure_cloud   | Contains Azure cloud resource IDs                                                                |
| aws_cloud     | Contains AWS cloud resource IDs                                                                  |
| gcp_cloud     | Contains GCP cloud resource IDs                                                                  |
| container     | Contains container IDs                                                                           |
| hostname      | Contains hostnames, not the FQDN, so only the first part (before `.`) will be used for matching. |
| user_identity | Contains user identities; f.e. username and/or e-mail                                            |
| ipv4          | IPv4 address or CIDR i.e. `10.1.2.0/24`                                                          |
| ipv6          | IPv6 address or CIDR i.e. `2a02:fe9:692:2812/64`                                                 |