package device42

import (
	"encoding/json"
	"testing"
)

func TestApiPduUnmarshal(t *testing.T) {
	// GET /1.0/pdus/{id}/ returns the PDU object directly (not wrapped in {"pdus": [...]}).
	body := `{"pdu_serial_no":"SN-001","name":"tf_ft_pdu","pdu_model":{"name":"AP7800","pdu_model_id":19},"pdu_id":248,"notes":"a note","category":"test","gateway_pdu_id":"7","custom_fields":[{"key":"k","notes":null,"value":"v"}]}`

	var r apiPdu
	if err := json.Unmarshal([]byte(body), &r); err != nil {
		t.Fatalf("unmarshal pdu: %s", err)
	}

	if r.PduId != 248 {
		t.Errorf("PduId = %d, want 248", r.PduId)
	}
	if r.Name != "tf_ft_pdu" {
		t.Errorf("Name = %q, want tf_ft_pdu", r.Name)
	}
	if r.PduSerialNo != "SN-001" {
		t.Errorf("PduSerialNo = %q, want SN-001", r.PduSerialNo)
	}
	if r.Notes != "a note" {
		t.Errorf("Notes = %q, want a note", r.Notes)
	}
	if r.GatewayPduId != "7" {
		t.Errorf("GatewayPduId = %q, want 7", r.GatewayPduId)
	}
	if r.PduModel == nil {
		t.Fatalf("PduModel is nil, want populated")
	}
	var model apiPduModel
	if err := json.Unmarshal(r.PduModel, &model); err != nil {
		t.Fatalf("unmarshal nested pdu_model: %s", err)
	}
	if model.Name != "AP7800" || model.PduModelId != 19 {
		t.Errorf("PduModel = %+v, want name=AP7800 id=19", model)
	}
}

func TestApiPduUnmarshalNoModel(t *testing.T) {
	// a PDU without a linked model returns an empty string for pdu_model, not null.
	body := `{"pdu_serial_no":"","name":"plain","pdu_model":"","pdu_id":1,"notes":"","category":"","custom_fields":[]}`

	var r apiPdu
	if err := json.Unmarshal([]byte(body), &r); err != nil {
		t.Fatalf("unmarshal pdu: %s", err)
	}
	if string(r.PduModel) != `""` {
		t.Errorf("PduModel = %q, want empty-string marker", string(r.PduModel))
	}
}
