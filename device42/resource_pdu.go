package device42

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiPdu struct {
	PduId        int64           `json:"pdu_id"`
	Name         string          `json:"name"`
	PduSerialNo  string          `json:"pdu_serial_no"`
	PduModel     json.RawMessage `json:"pdu_model"`
	Notes        string          `json:"notes"`
	Category     string          `json:"category"`
	CustomFields []customField   `json:"custom_fields"`
	GatewayPduId string          `json:"gateway_pdu_id"`
}

func resourceD42Pdu() *schema.Resource {
	return &schema.Resource{
		Description: "device42_pdu can be use to manage physical PDUs attached to a PDU model.",
		Create:      resourceDevice42PduCreate,
		Read:        resourceDevice42PduRead,
		Update:      resourceDevice42PduUpdate,
		Delete:      resourceDevice42PduDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pdu_serial_no": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"pdu_model": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pdu_model_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_pdu_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"custom_fields": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceDevice42PduCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := map[string]string{
		"name": d.Get("name").(string),
	}

	if v, ok := d.GetOk("pdu_serial_no"); ok {
		formData["pdu_serial_no"] = v.(string)
	}
	if v, ok := d.GetOk("pdu_model"); ok {
		formData["pdu_model"] = v.(string)
	}
	if v, ok := d.GetOk("notes"); ok {
		formData["notes"] = v.(string)
	}
	if v, ok := d.GetOk("gateway_pdu_id"); ok {
		formData["gateway_pdu_id"] = v.(string)
	}

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/pdus/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42PduCreate - API returned code %d: %v", r.Code, r.Msg)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42PduCreate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42PduRead(d, m)
}

func resourceDevice42PduRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42PduRead - Starting reading using API for id %s", d.Id())

	var resp apiPdu
	apiResp, err := client.R().
		SetResult(&resp).
		Get(fmt.Sprintf("/1.0/pdus/%s/", d.Id()))
	if err != nil {
		log.Printf("[WARN] No pdu found for id %s", d.Id())
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() == 404 {
		log.Printf("[WARN] No pdu found for id %s", d.Id())
		d.SetId("")
		return nil
	}
	if apiResp.StatusCode() >= 400 {
		return fmt.Errorf("resourceDevice42PduRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	p := resp
	d.Set("name", p.Name)
	d.Set("pdu_serial_no", p.PduSerialNo)
	if len(p.PduModel) > 0 && string(p.PduModel) != `""` {
		var model apiPduModel
		if err := json.Unmarshal(p.PduModel, &model); err != nil {
			return fmt.Errorf("resourceDevice42PduRead - malformed pdu_model: %v", err)
		}
		d.Set("pdu_model", model.Name)
		d.Set("pdu_model_id", strconv.FormatInt(model.PduModelId, 10))
	}
	d.Set("notes", p.Notes)
	d.Set("gateway_pdu_id", p.GatewayPduId)
	d.Set("custom_fields", flattenCustomFields(p.CustomFields))

	return nil
}

func resourceDevice42PduUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	payload := map[string]string{
		"id":   d.Id(),
		"name": d.Get("name").(string),
	}

	if v, ok := d.GetOk("pdu_serial_no"); ok {
		payload["pdu_serial_no"] = v.(string)
	}
	if v, ok := d.GetOk("pdu_model"); ok {
		payload["pdu_model"] = v.(string)
	}
	if v, ok := d.GetOk("notes"); ok {
		payload["notes"] = v.(string)
	}

	resp, err := client.R().
		SetResult(apiResponse{}).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Put("/1.0/pdus/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42PduUpdate - API returned code %d: %v", r.Code, r.Msg)
	}
	log.Printf("[DEBUG] Result: %#v", r)

	return resourceDevice42PduRead(d, m)
}

func resourceDevice42PduDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting pdu id %s (name: %s)", d.Id(), d.Get("name"))

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(fmt.Sprintf("/1.0/pdus/%s/", d.Id()))

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42PduDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}
