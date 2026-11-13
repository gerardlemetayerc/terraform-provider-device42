package device42

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// apiSoftwareLicenseKeyResponse is the shape of a single software license key
// returned by GET /1.0/software/license_keys/. The endpoint returns ALL keys in
// one response, so reads are filtered client-side by id.
type apiSoftwareLicenseKeyResponse struct {
	ID           int64  `json:"id"`
	Key          string `json:"key"`
	SoftwareName string `json:"software_name"`
	SoftwareID   int64  `json:"software_id"`
	Count        int64  `json:"count"`
	Notes        string `json:"notes"`
}

// apiSoftwareLicenseKeyListResponse is the wrapper returned by
// GET /1.0/software/license_keys/.
type apiSoftwareLicenseKeyListResponse struct {
	TotalCount          int64                           `json:"total_count"`
	SoftwareLicenseKeys []apiSoftwareLicenseKeyResponse `json:"software_license_keys"`
}

func resourceD42SoftwareLicenseKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceDevice42SoftwareLicenseKeyCreate,
		Read:   resourceDevice42SoftwareLicenseKeyRead,
		Update: resourceDevice42SoftwareLicenseKeyUpdate,
		Delete: resourceDevice42SoftwareLicenseKeyDelete,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "The license key value. Used as the upsert key by the Device42 API.",
			},
			"software_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the software component this license key belongs to.",
			},
			"software_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the software component (derived by the API).",
			},
			"license_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of licenses covered by this key.",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Any additional notes.",
			},
			"license_key_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the license key.",
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				// The import key is the Device42 license key id. Validate it is
				// numeric so a malformed key fails fast; the Read then populates the
				// remaining fields by filtering the list.
				id, err := strconv.Atoi(d.Id())
				if err != nil {
					return nil, fmt.Errorf("expected license key id to be a numeric id, got %q", d.Id())
				}
				d.SetId(strconv.Itoa(id))
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceDevice42SoftwareLicenseKeyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	formData := map[string]string{}
	formData["key"] = d.Get("key").(string)
	formData["software_id"] = strconv.Itoa(d.Get("software_id").(int))

	if v := d.Get("license_count").(int); v > 0 {
		formData["count"] = strconv.Itoa(v)
	}
	if v := d.Get("notes").(string); v != "" {
		formData["notes"] = v
	}

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/software/license_keys/")
	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42SoftwareLicenseKeyCreate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42SoftwareLicenseKeyCreate - Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42SoftwareLicenseKeyCreate - unable to extract license key id from response: %s", err)
	}

	d.SetId(strconv.Itoa(id))
	d.Set("license_key_id", id)

	return resourceDevice42SoftwareLicenseKeyRead(d, m)
}

func resourceDevice42SoftwareLicenseKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42SoftwareLicenseKeyRead - starting read for license key id %s", d.Id())

	var resp apiSoftwareLicenseKeyListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/software/license_keys/")
	if err != nil {
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42SoftwareLicenseKeyRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("resourceDevice42SoftwareLicenseKeyRead - invalid license key id %q: %s", d.Id(), err)
	}

	for _, k := range resp.SoftwareLicenseKeys {
		if int(k.ID) == id {
			d.Set("key", k.Key)
			d.Set("software_id", k.SoftwareID)
			d.Set("software_name", k.SoftwareName)
			d.Set("license_count", k.Count)
			d.Set("notes", k.Notes)
			d.Set("license_key_id", k.ID)
			return nil
		}
	}

	log.Printf("[WARN] resourceDevice42SoftwareLicenseKeyRead - no license key found: %s", d.Id())
	d.SetId("")
	return nil
}

func resourceDevice42SoftwareLicenseKeyUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	formData := map[string]string{}
	formData["id"] = strconv.Itoa(d.Get("license_key_id").(int))
	formData["key"] = d.Get("key").(string)
	formData["software_id"] = strconv.Itoa(d.Get("software_id").(int))

	if d.HasChange("license_count") {
		formData["count"] = strconv.Itoa(d.Get("license_count").(int))
	}
	if d.HasChange("notes") {
		formData["notes"] = d.Get("notes").(string)
	}

	log.Printf("[DEBUG] resourceDevice42SoftwareLicenseKeyUpdate - updating license key ID %s", d.Id())

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/software/license_keys/")
	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42SoftwareLicenseKeyUpdate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42SoftwareLicenseKeyUpdate - Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42SoftwareLicenseKeyUpdate - unable to extract license key id from response: %s", err)
	}

	d.SetId(strconv.Itoa(id))
	d.Set("license_key_id", id)

	return resourceDevice42SoftwareLicenseKeyRead(d, m)
}

func resourceDevice42SoftwareLicenseKeyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting software license key ID %s", d.Id())

	url := fmt.Sprintf("/1.0/software/license_keys/%s/", d.Id())

	resp, err := client.R().
		SetResult(apiDeleteResponse{}).
		Delete(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42SoftwareLicenseKeyDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiDeleteResponse)
	log.Printf("[DEBUG] resourceDevice42SoftwareLicenseKeyDelete - Result: %#v", r)
	return nil
}
