resource "device42_device" "myNewVM" {
  name                 = "myNewVMName"
  type                 = "virtual"
  archive_on_destroy   = true
  asset_no             = "ASSET-0001"
  building             = "Main Data Center"
  customer             = "Acme Corp"
}