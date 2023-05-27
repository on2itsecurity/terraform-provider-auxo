resource "auxo_transactionflow" "tf_ps_ad" {
  protectsurface                 = auxo_protectsurface.ps_ad.id
  incoming_protectsurfaces_allow = [auxo_protectsurface.ps_mail.id]
  incoming_protectsurfaces_block = [auxo_protectsurface.ps_guests.id]
  outgoing_protectsurfaces_allow = [auxo_protectsurface.ps_mail.id]
  outgoing_protectsurfaces_block = [auxo_protectsurface.ps_guests.id]
}