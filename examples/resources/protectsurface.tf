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
  allow_flows_from_outside = false
  allow_flows_to_outside   = false
}
