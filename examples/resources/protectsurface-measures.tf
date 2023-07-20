 data "auxo_protectsurface" "ps_mail" {
  name = "Mail"
}

data "auxo_contact" "rob" {
  email = "rob.maas+tst@on2it.net"
}
 
resource auxo_measure ps_mail {
  protectsurface = data.auxo_protectsurface.ps_mail.id
  measures = {
    flows-segmentation = {
      assigned                        = true
      assigned_by                     = data.auxo_contact.rob.email
      implemented                     = true
      implemented_by                  = data.auxo_contact.rob.email
      evidenced                       = false
      evidenced_by                    = data.auxo_contact.rob.email
      risk_no_implementation_accepted = false
      risk_acceptance_by              = data.auxo_contact.rob.email
      risk_no_evidence_accepted       = true
      risk_accepted_comment           = "This is a test"

    },
    encryption-at-rest = {
      assigned       = true
      assigned_by    = data.auxo_contact.rob.email
      implemented    = true
      implemented_by = data.auxo_contact.rob.email
      evidenced      = false
      evidenced_by   = data.auxo_contact.rob.email
    },
    encryption-in-transit = {
      assigned       = true
      assigned_by    = data.auxo_contact.rob.email
      implemented    = true
      implemented_by = data.auxo_contact.rob.email
      evidenced      = false
      evidenced_by   = data.auxo_contact.rob.email
    }
  }
}

