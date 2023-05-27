resource "auxo_state" "ps_ad-loc_zaltbommel-ipv4" {
  content_type      = "ipv4"
  description       = "IPv4 allocations of AD servers"
  location_id       = auxo_location.loc_zaltbommel.id
  protectsurface_id = auxo_protectsurface.ps_ad.id
  content           = ["10.0.42.10", "10.0.42.11", "10.0.42.12"]
}
