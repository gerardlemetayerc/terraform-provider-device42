resource "device42_vrf_group" "myNewVrfGroup" {
  name        = "myNewVrfGroup"
  description = "my new VRF group"
  tags        = "prod,core"
}