resource "device42_rack" "myNewRack" {
    name          = "myNewRack"
    size          = 42
    building      = "myNewBuilding"
    room          = "myNewRoom"
    manufacturer  = "APC"
    tags          = "Prod,East"
    custom_fields = {
        "rack_type" = "server"
    }
}