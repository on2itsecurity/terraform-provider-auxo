// The used provider
terraform {
  required_providers {
    auxo = {
      version = "0.0.1"
      source  = "on2itsecurity/auxo"
    }
  }
}

// Represents protect-surface "Active Directory"
resource "auxo_protectsurface" "ps_ad" {
  name                = "Active Directory"
  description         = "Active Directory for employees"
  relevance           = 90
  in_control_boundary = true
  in_zero_trust_focus = true
  confidentiality     = 3
  availability        = 2
  integrity           = 3
  main_contact        = "37904"
  security_contact    = "37904"
  data_tags           = ["PII"]
  compliance_tags     = ["GDPR"]
  customer_labels = {
    owner      = "Rob Maas"
    env        = "Production"
    os         = "Windows"
    created-by = "Terraform"
  }
  soc_tags       = ["active-directory", "windows"]
}

// Represents protect-surface "Mail"
resource "auxo_protectsurface" "ps_mail" {
  name                = "Mail"
  description         = "Mail and all related componentes e.g. Spamfilter"
  relevance           = 80
  confidentiality     = 2
  availability        = 2
  integrity           = 3
  in_control_boundary = true
  in_zero_trust_focus = true
  main_contact        = "37904"
  security_contact    = "37904"
  data_tags           = ["PII", "PCI"]
  compliance_tags     = ["PCI"]
  soc_tags            = ["postfix", "linux"]
  customer_labels = {
    owner      = "Rob Maas"
    env        = "Production"
    os         = "Linux"
    created-by = "Terraform"
  }
}

// Represents protect-surface "Guests"
resource "auxo_protectsurface" "ps_guests" {
  name                = "Guests"
  description         = "Guest network"
  relevance           = 10
  in_control_boundary = true
  in_zero_trust_focus = false
}

// Represents location (this is used in states to specify where resources live)
resource "auxo_location" "loc_beuningen" {
  name      = "Datacenter Beuningen"
  latitude  = 51.8557833
  longitude = 5.7490162
}

// Represents location (this is used in states to specify where resources live)
resource "auxo_location" "loc_zaltbommel" {
  name      = "Datacenter Zaltbommel"
  latitude  = 51.7983645
  longitude = 5.2548381
}


// Represents a state, this can be seen as the glue between protect-surface, location and the type of resources it represent
// A protect surface can have multipe states attached to it
resource "auxo_state" "ps_ad-loc_beunigen-ipv4" {
  content_type      = "static_ipv4"
  description       = "Static IPv4 allocations of AD servers"
  location_id       = auxo_location.loc_beuningen.id
  protectsurface_id = auxo_protectsurface.ps_ad.id
  content           = ["10.42.0.10", "10.42.0.11", "10.42.0.12"]
}

resource "auxo_state" "ps_ad-loc_zaltbommel-ipv4" {
  content_type      = "static_ipv4"
  description       = "Static IPv4 allocations of AD servers"
  location_id       = auxo_location.loc_zaltbommel.id
  protectsurface_id = auxo_protectsurface.ps_ad.id
  content           = ["10.0.42.10", "10.0.42.11", "10.0.42.12"]
}

resource "auxo_state" "ps_guests-loc_zaltbommel-ipv4" {
  content_type      = "static_ipv4"
  description       = "Static IPv4 subnet for guests"
  location_id       = auxo_location.loc_zaltbommel.id
  protectsurface_id = auxo_protectsurface.ps_guests.id
  content           = ["192.168.42.0/24"]
}
