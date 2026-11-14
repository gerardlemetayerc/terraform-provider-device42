package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiPowerCircuitResponse struct {
	PowerCircuitId int64   `json:"power_circuit_id"`
	BreakerpanelId int64   `json:"breakerpanel_id"`
	Breakerpanel   string  `json:"breakerpanel"`
	Number         string  `json:"number"`
	Label          string  `json:"label"`
	Amps           float64 `json:"amps"`
	Voltage        float64 `json:"voltage"`
	Connection     string  `json:"connection"`
}

func resourceD42PowerCircuit() *schema.Resource {
	return &schema.Resource{
		Description: "device42_power_circuit can be use to manage power circuits attached to a breaker panel.",
		Create:      resourceDevice42PowerCircuitCreate,
		Read:        resourceDevice42PowerCircuitRead,
		Update:      resourceDevice42PowerCircuitUpdate,
		Delete:      resourceDevice42PowerCircuitDelete,

		Schema: map[string]*schema.Schema{
			"breakerpanel_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"number": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"breakerpanel": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"amps": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"voltage": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"connection_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDevice42PowerCircuitCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := powerCircuitFormData(d, "")

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/powercircuits/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42PowerCircuitCreate - API returned code %d: %v", r.Code, r.Msg)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42PowerCircuitCreate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42PowerCircuitRead(d, m)
}

func resourceDevice42PowerCircuitRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42PowerCircuitRead - Starting reading using API for id %s", d.Id())
	resp, err := client.R().
		SetResult(apiPowerCircuitResponse{}).
		Get(fmt.Sprintf("/1.0/powercircuits/%s/", d.Id()))

	if err != nil {
		log.Printf("[WARN] No power circuit found for id %s", d.Id())
		d.SetId("")
		return err
	}
	if resp.StatusCode() == 404 {
		log.Printf("[WARN] No power circuit found for id %s", d.Id())
		d.SetId("")
		return nil
	}
	if resp.StatusCode() >= 400 {
		return fmt.Errorf("resourceDevice42PowerCircuitRead - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiPowerCircuitResponse)
	d.Set("breakerpanel_id", strconv.FormatInt(r.BreakerpanelId, 10))
	d.Set("number", r.Number)
	d.Set("breakerpanel", r.Breakerpanel)
	d.Set("label", r.Label)
	d.Set("amps", int64(r.Amps))
	d.Set("voltage", int64(r.Voltage))
	d.Set("connection_name", r.Connection)

	return nil
}

func resourceDevice42PowerCircuitUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := powerCircuitFormData(d, d.Id())

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/powercircuits/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42PowerCircuitUpdate - API returned code %d: %v", r.Code, r.Msg)
	}
	log.Printf("[DEBUG] Result: %#v", r)

	return resourceDevice42PowerCircuitRead(d, m)
}

func resourceDevice42PowerCircuitDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting power circuit id %s (number: %s)", d.Id(), d.Get("number"))

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(fmt.Sprintf("/1.0/powercircuits/%s/", d.Id()))

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42PowerCircuitDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

func powerCircuitFormData(d *schema.ResourceData, id string) map[string]string {
	formData := map[string]string{
		"breakerpanel_id": d.Get("breakerpanel_id").(string),
		"number":          d.Get("number").(string),
	}

	if id != "" {
		formData["id"] = id
	}

	if v, ok := d.GetOk("label"); ok {
		formData["label"] = v.(string)
	}
	if v, ok := d.GetOk("amps"); ok {
		formData["amps"] = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("voltage"); ok {
		formData["voltage"] = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("connection_name"); ok {
		formData["connection"] = v.(string)
	}

	return formData
}
