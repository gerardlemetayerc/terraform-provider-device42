resource "device42_pdu" "myNewPdu" {
  name          = "myNewPdu"
  pdu_serial_no = "SN-001"
  pdu_model     = device42_pdu_model.myNewPduModel.name
  notes         = "my new PDU"
}