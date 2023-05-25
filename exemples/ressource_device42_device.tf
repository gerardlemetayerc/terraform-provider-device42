resource "device42_device" "myNewVM" {
  name                 = "myNewVMName"
  type                 = "virtual"
  archive_on_destroy   = true
}