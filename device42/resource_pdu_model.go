package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiPduModel struct {
	PduModelId                  int64   `json:"pdu_model_id"`
	Name                        string  `json:"name"`
	Type                        string  `json:"type"`
	Manufacturer                string  `json:"manufacturer"`
	Description                 string  `json:"description"`
	Notes                       string  `json:"notes"`
	PortCount                   int64   `json:"port_count"`
	Depth                       string  `json:"depth"`
	Size                        float64 `json:"size"`
	Width                       float64 `json:"width"`
	Height                      float64 `json:"height"`
	SequentialNumberingForPorts bool    `json:"sequential_numbering_for_ports"`
}

type apiPduModelListResponse struct {
	PduModels []apiPduModel `json:"pdu_models"`
}

func resourceD42PduModel() *schema.Resource {
	return &schema.Resource{
		Description: "device42_pdu_model can be use to manage PDU models (the model/type of a physical PDU).",
		Create:      resourceDevice42PduModelCreate,
		Read:        resourceDevice42PduModelRead,
		Update:      resourceDevice42PduModelUpdate,
		Delete:      resourceDevice42PduModelDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"manufacturer": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"depth": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"width": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"height": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceDevice42PduModelCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := pduModelFormData(d, "")

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/pdu_models/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42PduModelCreate - API returned code %d: %v", r.Code, r.Msg)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42PduModelCreate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42PduModelRead(d, m)
}

func resourceDevice42PduModelRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42PduModelRead - Starting reading using API for id %s", d.Id())

	var resp apiPduModelListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/pdu_models/")
	if err != nil {
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42PduModelRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("resourceDevice42PduModelRead - invalid id %q: %s", d.Id(), err)
	}

	for _, m := range resp.PduModels {
		if int(m.PduModelId) == id {
			d.Set("name", m.Name)
			d.Set("manufacturer", m.Manufacturer)
			d.Set("description", m.Description)
			d.Set("notes", m.Notes)
			d.Set("depth", m.Depth)
			d.Set("size", int64(m.Size))
			d.Set("width", int64(m.Width))
			d.Set("height", int64(m.Height))
			return nil
		}
	}

	log.Printf("[WARN] resourceDevice42PduModelRead - no pdu model found: %s", d.Id())
	d.SetId("")
	return nil
}

func resourceDevice42PduModelUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := pduModelFormData(d, d.Id())

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/pdu_models/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42PduModelUpdate - API returned code %d: %v", r.Code, r.Msg)
	}
	log.Printf("[DEBUG] Result: %#v", r)

	return resourceDevice42PduModelRead(d, m)
}

func resourceDevice42PduModelDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting pdu model id %s (name: %s)", d.Id(), d.Get("name"))

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(fmt.Sprintf("/1.0/pdu_models/%s/", d.Id()))

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42PduModelDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

func pduModelFormData(d *schema.ResourceData, id string) map[string]string {
	formData := map[string]string{
		"name": d.Get("name").(string),
	}

	if id != "" {
		formData["id"] = id
	}

	if v, ok := d.GetOk("manufacturer"); ok {
		formData["manufacturer"] = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		formData["description"] = v.(string)
	}
	if v, ok := d.GetOk("notes"); ok {
		formData["notes"] = v.(string)
	}
	if v, ok := d.GetOk("depth"); ok {
		formData["depth"] = v.(string)
	}
	if v, ok := d.GetOk("size"); ok {
		formData["size"] = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("width"); ok {
		formData["width"] = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("height"); ok {
		formData["height"] = strconv.Itoa(v.(int))
	}

	return formData
}
