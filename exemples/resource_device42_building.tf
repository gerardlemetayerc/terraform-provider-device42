resource "device42_building" "myNewBuilding" {
    name           = "myNewBuilding"
    address        = "100 Main Street, New York"
    contact_name   = "John Doe"
    contact_phone  = "+1 555 000 000"
    notes          = "Main Data Center building"
    tags           = "Prod,East"
    custom_fields  = {
        "building_type" = "data_center"
    }
}