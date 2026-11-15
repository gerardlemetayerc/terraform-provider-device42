resource "device42_auto_discovery_pingsweep" "example" {
  name              = "myPingSweepJob"
  networks          = "192.168.0.0/24"
  create_new_subnet = "yes"
  schedule_days     = "0,1,2,3,4,5,6"
}