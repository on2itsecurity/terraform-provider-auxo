// The used provider
terraform {
  required_providers {
    auxo = {
      source = "on2itsecurity/auxo"
    }
  }
}

provider "auxo" {
  name = "tailspin"
}

// Get the contact based on the email address
data "auxo_contact" "rob" {
  email = "rob.maas+tst@on2it.net"
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
  main_contact        = data.auxo_contact.rob.id
  security_contact    = data.auxo_contact.rob.id
  data_tags           = ["PII"]
  compliance_tags     = ["GDPR"]
  customer_labels = {
    owner      = "Rob Maas"
    env        = "Production"
    os         = "Windows"
    created-by = "Terraform"
  }
  soc_tags                 = ["active-directory", "windows"]
  allow_flows_from_outside = false
  allow_flows_to_outside   = false

  // Represents the segmentation measure for this protect-surface
  measures = {
    flows-segmentation = {
      assigned       = true
      assigned_by    = data.auxo_contact.rob.email
      implemented    = true
      implemented_by = data.auxo_contact.rob.email
      evidenced      = false
      evidenced_by   = data.auxo_contact.rob.email
    } 
  }
}

// Represents transactionflows related to protect surface "Active Directory"
resource "auxo_transactionflow" "tf_ps_ad" {
  protectsurface                 = auxo_protectsurface.ps_ad.id
  incoming_protectsurfaces_allow = [auxo_protectsurface.ps_mail.id]
  incoming_protectsurfaces_block = [auxo_protectsurface.ps_guests.id]
  outgoing_protectsurfaces_allow = [auxo_protectsurface.ps_mail.id]
  outgoing_protectsurfaces_block = [auxo_protectsurface.ps_guests.id]
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
  main_contact        = data.auxo_contact.rob.id
  security_contact    = data.auxo_contact.rob.id
  data_tags           = ["PII", "PCI"]
  compliance_tags     = ["PCI"]
  soc_tags            = ["postfix", "linux"]
  customer_labels = {
    owner      = "Rob Maas"
    env        = "Production"
    os         = "Linux"
    created-by = "Terraform"
  }
  allow_flows_from_outside = true
  allow_flows_to_outside   = true
}

// Represents transactionflows related to protect surface "Mail"
resource "auxo_transactionflow" "tf_ps_mail" {
  protectsurface                 = auxo_protectsurface.ps_mail.id
  incoming_protectsurfaces_allow = [auxo_protectsurface.ps_ad.id]
  incoming_protectsurfaces_block = [auxo_protectsurface.ps_guests.id]
  outgoing_protectsurfaces_allow = [auxo_protectsurface.ps_ad.id]
  outgoing_protectsurfaces_block = [auxo_protectsurface.ps_guests.id]
}

// Represents protect-surface "Guests"
resource "auxo_protectsurface" "ps_guests" {
  name                     = "Guests"
  description              = "Guest network"
  relevance                = 10
  in_control_boundary      = true
  in_zero_trust_focus      = false
  allow_flows_to_outside   = true
  allow_flows_from_outside = false
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
  content_type      = "ipv4"
  description       = "Static IPv4 allocations of AD servers"
  location_id       = auxo_location.loc_beuningen.id
  protectsurface_id = auxo_protectsurface.ps_ad.id
  content           = ["10.42.0.10", "10.42.0.11", "10.42.0.12"]
}

resource "auxo_state" "ps_ad-loc_zaltbommel-ipv4" {
  content_type      = "ipv4"
  description       = "Static IPv4 allocations of AD servers"
  location_id       = auxo_location.loc_zaltbommel.id
  protectsurface_id = auxo_protectsurface.ps_ad.id
  content           = ["10.0.42.10", "10.0.42.11", "10.0.42.12"]
}

resource "auxo_state" "ps_guests-loc_zaltbommel-ipv4" {
  content_type      = "ipv4"
  description       = "Static IPv4 subnet for guests"
  location_id       = auxo_location.loc_zaltbommel.id
  protectsurface_id = auxo_protectsurface.ps_guests.id
  content           = ["192.168.42.0/24"]
}
