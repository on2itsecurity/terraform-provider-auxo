resource "auxo_transactionflow" "tf_ps_a" {
  protectsurface                 = auxo_protectsurface.ps_a.id
  outgoing_protectsurfaces_allow = [auxo_protectsurface.ps_b.id]
}

resource "auxo_transactionflow" "tf_ps_b" {
  protectsurface                 = auxo_protectsurface.ps_b.id
  incoming_protectsurfaces_allow = [auxo_protectsurface.ps_a.id]
}