package device42

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// apiVendorReadResponse is the shape of a single vendor object returned by
// GET /1.0/vendors/. GET returns ALL vendors (no single-id endpoint), so reads
// are filtered client-side by vendor_id.
type apiVendorReadResponse struct {
	Name         string        `json:"name"`
	VendorID     int64         `json:"vendor_id"`
	HomePage     string        `json:"home_page"`
	Escalation1  string        `json:"escalation_1"`
	Escalation2  string        `json:"escalation_2"`
	Notes        string        `json:"notes"`
	AccountNo    string        `json:"account_no"`
	Aliases      string        `json:"aliases"`
	ContactInfo  string        `json:"contact_info"`
	Tags         []interface{} `json:"tags"`
	CustomFields []customField `json:"custom_fields"`
}

// apiVendorsResponse is the wrapper returned by GET /1.0/vendors/.
type apiVendorsResponse struct {
	Vendors []apiVendorReadResponse `json:"vendors"`
}

func resourceD42Vendor() *schema.Resource {
	return &schema.Resource{
		Create: resourceDevice42VendorCreate,
		Read:   resourceDevice42VendorRead,
		Update: resourceDevice42VendorUpdate,
		Delete: resourceDevice42VendorDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the vendor. Used as the upsert key by the Device42 API.",
			},
			"vendor_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the vendor.",
			},
			"home_page": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Home page URL of the vendor.",
			},
			"contact_info": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Contact information for the vendor.",
			},
			"account_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Account number associated with the vendor.",
			},
			"aliases": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Aliases of the vendor.",
			},
			"escalation_1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "First escalation contact.",
			},
			"escalation_2": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Second escalation contact.",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Any additional notes.",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Add or update tags to a vendor (comma separated).",
			},
			"custom_fields": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				Description:      "Any custom fields that will be used in device42.",
				DiffSuppressFunc: suppressCustomFieldsDiffs,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				// The import key is the Device42 vendor_id. Validate it is numeric
				// so a malformed key fails fast; resourceDevice42VendorRead then
				// populates the remaining fields by filtering the full list.
				id, err := strconv.Atoi(d.Id())
				if err != nil {
					return nil, fmt.Errorf("expected vendor id to be a numeric id, got %q", d.Id())
				}
				d.SetId(strconv.Itoa(id))
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceDevice42VendorCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	name := d.Get("name").(string)
	formData := map[string]string{}
	formData["name"] = name

	if v := d.Get("home_page").(string); v != "" {
		formData["home_page"] = v
	}
	if v := d.Get("contact_info").(string); v != "" {
		formData["contact_info"] = v
	}
	if v := d.Get("account_no").(string); v != "" {
		formData["account_no"] = v
	}
	if v := d.Get("aliases").(string); v != "" {
		formData["aliases"] = v
	}
	if v := d.Get("escalation_1").(string); v != "" {
		formData["escalation_1"] = v
	}
	if v := d.Get("escalation_2").(string); v != "" {
		formData["escalation_2"] = v
	}
	if v := d.Get("notes").(string); v != "" {
		formData["notes"] = v
	}
	if v := d.Get("tags").(string); v != "" {
		formData["tags"] = v
	}

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/vendors/")
	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42VendorCreate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42VendorCreate - Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42VendorCreate - unable to extract vendor id from response: %s", err)
	}

	d.SetId(strconv.Itoa(id))
	d.Set("vendor_id", id)

	if d.Get("custom_fields") != nil {
		fields := d.Get("custom_fields").(map[string]interface{})
		bulkFields := []string{}
		for k, v := range fields {
			bulkFields = append(bulkFields, fmt.Sprintf("%v:%v", k, v))
		}

		log.Printf("[DEBUG] resourceDevice42VendorCreate - Pushing custom_fields %s", bulkFields)
		resp, err := client.R().
			SetFormData(map[string]string{
				"name":        name,
				"bulk_fields": strings.Join(bulkFields, ","),
			}).
			SetResult(apiResponse{}).
			Put("/1.0/custom_fields/vendor/")
		if err != nil {
			return err
		}

		r := resp.Result().(*apiResponse)
		if r.Code != 0 {
			return fmt.Errorf("resourceDevice42VendorCreate - api returned code %d", r.Code)
		}
	}

	return resourceDevice42VendorRead(d, m)
}

func resourceDevice42VendorRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42VendorRead - starting read for vendor id %s", d.Id())

	var resp apiVendorsResponse
	_, err := client.R().
		SetResult(&resp).
		Get("/1.0/vendors/")
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("resourceDevice42VendorRead - invalid vendor id %q: %s", d.Id(), err)
	}

	for _, v := range resp.Vendors {
		if int(v.VendorID) == id {
			d.Set("name", v.Name)
			d.Set("vendor_id", v.VendorID)
			d.Set("home_page", v.HomePage)
			d.Set("contact_info", v.ContactInfo)
			d.Set("account_no", v.AccountNo)
			d.Set("aliases", v.Aliases)
			d.Set("escalation_1", v.Escalation1)
			d.Set("escalation_2", v.Escalation2)
			d.Set("notes", v.Notes)
			d.Set("tags", joinBuildingList(v.Tags))
			d.Set("custom_fields", flattenCustomFields(v.CustomFields))
			return nil
		}
	}

	log.Printf("[WARN] resourceDevice42VendorRead - no vendor found: %s", d.Id())
	d.SetId("")
	return nil
}

func resourceDevice42VendorUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	name := d.Get("name").(string)
	formData := map[string]string{}
	formData["name"] = name

	if d.HasChange("home_page") {
		formData["home_page"] = d.Get("home_page").(string)
	}
	if d.HasChange("contact_info") {
		formData["contact_info"] = d.Get("contact_info").(string)
	}
	if d.HasChange("account_no") {
		formData["account_no"] = d.Get("account_no").(string)
	}
	if d.HasChange("aliases") {
		formData["aliases"] = d.Get("aliases").(string)
	}
	if d.HasChange("escalation_1") {
		formData["escalation_1"] = d.Get("escalation_1").(string)
	}
	if d.HasChange("escalation_2") {
		formData["escalation_2"] = d.Get("escalation_2").(string)
	}
	if d.HasChange("notes") {
		formData["notes"] = d.Get("notes").(string)
	}
	if d.HasChange("tags") {
		formData["tags"] = d.Get("tags").(string)
	}

	log.Printf("[DEBUG] resourceDevice42VendorUpdate - Pushing new informations : %#v", formData)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/vendors/")
	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42VendorUpdate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42VendorUpdate - Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42VendorUpdate - unable to extract vendor id from response: %s", err)
	}

	d.SetId(strconv.Itoa(id))
	d.Set("vendor_id", id)

	if d.HasChange("custom_fields") {
		updateList := setCustomFields(d)
		for k, v := range updateList {
			log.Printf("[DEBUG] resourceDevice42VendorUpdate - Pushing custom field %s : %v", k, v)
			resp, err := client.R().
				SetFormData(map[string]string{
					"name":  name,
					"key":   k,
					"value": v.(string),
				}).
				SetResult(apiResponse{}).
				Put("/1.0/custom_fields/vendor/")
			if err != nil {
				return err
			}

			r := resp.Result().(*apiResponse)
			if resp.IsError() || r.Code != 0 {
				return fmt.Errorf("resourceDevice42VendorUpdate - custom field %q update failed (HTTP %d, API code %d)", k, resp.StatusCode(), r.Code)
			}
			log.Printf("[DEBUG] resourceDevice42VendorUpdate - Result: %#v", r)
		}
	}

	return resourceDevice42VendorRead(d, m)
}

func resourceDevice42VendorDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting vendor %s (ID: %s)", d.Get("name"), d.Id())

	url := fmt.Sprintf("/1.0/vendors/%s/", d.Id())

	resp, err := client.R().
		SetResult(apiDeleteResponse{}).
		Delete(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42VendorDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiDeleteResponse)
	log.Printf("[DEBUG] resourceDevice42VendorDelete - Result: %#v", r)
	return nil
}
