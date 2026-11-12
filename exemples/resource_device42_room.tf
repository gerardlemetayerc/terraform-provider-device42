resource "device42_room" "myNewRoom" {
    name          = "myNewRoom"
    building_id   = 1
    notes         = "Main server room"
    tags          = "Prod,East"
}