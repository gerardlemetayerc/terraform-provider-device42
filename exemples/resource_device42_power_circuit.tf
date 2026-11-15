resource "device42_power_circuit" "myNewPowerCircuit" {
  breakerpanel_id = "30"
  number          = "1"
  label           = "myNewPowerCircuit"
  amps            = 20
  voltage         = 120
  connection_name = "main"
}