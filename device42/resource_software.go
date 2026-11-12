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

// apiSoftwareReadResponse is the shape of a single software component object
// returned by GET /1.0/software/. The endpoint returns ALL components in one
// response (no pagination), so reads are filtered client-side by id.
type apiSoftwareReadResponse struct {
	ID                       int64         `json:"id"`
	Name                     string        `json:"name"`
	Vendor                   string        `json:"vendor"`
	Description              string        `json:"description"`
	Notes                    string        `json:"notes"`
	SoftwareType             string        `json:"software_type"`
	DeploymentType           string        `json:"deployment_type"`
	Category                 string        `json:"category"`
	Aliases                  string        `json:"aliases"`
	LicensingModel           string        `json:"licensing_model"`
	TrackLicensedCountByKeys string        `json:"track_licensed_count_by_keys"`
	LicensedCount            float64       `json:"licensed_count"`
	DiscoveredCount          float64       `json:"discovered_count"`
	Tags                     []interface{} `json:"tags"`
	CustomFields             []customField `json:"custom_fields"`
}

// apiSoftwareListResponse is the wrapper returned by GET /1.0/software/.
type apiSoftwareListResponse struct {
	TotalCount int64                     `json:"total_count"`
	Software   []apiSoftwareReadResponse `json:"software"`
}

func resourceD42Software() *schema.Resource {
	return &schema.Resource{
		Create: resourceDevice42SoftwareCreate,
		Read:   resourceDevice42SoftwareRead,
		Update: resourceDevice42SoftwareUpdate,
		Delete: resourceDevice42SoftwareDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the software component. Used as the upsert key by the Device42 API.",
			},
			"software_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the software component.",
			},
			"vendor": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Vendor of the software component.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the software component.",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Any additional notes.",
			},
			"software_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the software component.",
			},
			"deployment_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Deployment type of the software component.",
			},
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Category of the software component.",
			},
			"aliases": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Aliases of the software component.",
			},
			"licensing_model": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Licensing model of the software component.",
			},
			"track_licensed_count_by_keys": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether licensed counts are tracked by license keys (yes/no).",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comma-separated tags to attach to the software component.",
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
				// The import key is the Device42 software id. Validate it is numeric
				// so a malformed key fails fast; resourceDevice42SoftwareRead then
				// populates the remaining fields by filtering the full list.
				id, err := strconv.Atoi(d.Id())
				if err != nil {
					return nil, fmt.Errorf("expected software id to be a numeric id, got %q", d.Id())
				}
				d.SetId(strconv.Itoa(id))
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceDevice42SoftwareCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	name := d.Get("name").(string)
	formData := map[string]string{}
	formData["name"] = name

	if v := d.Get("vendor").(string); v != "" {
		formData["vendor"] = v
	}
	if v := d.Get("description").(string); v != "" {
		formData["description"] = v
	}
	if v := d.Get("notes").(string); v != "" {
		formData["notes"] = v
	}
	if v := d.Get("software_type").(string); v != "" {
		formData["software_type"] = v
	}
	if v := d.Get("deployment_type").(string); v != "" {
		formData["deployment_type"] = v
	}
	if v := d.Get("category").(string); v != "" {
		formData["category"] = v
	}
	if v := d.Get("aliases").(string); v != "" {
		formData["aliases"] = v
	}
	if v := d.Get("licensing_model").(string); v != "" {
		formData["licensing_model"] = v
	}
	if v := d.Get("track_licensed_count_by_keys").(string); v != "" {
		formData["track_licensed_count_by_keys"] = v
	}
	if v := d.Get("tags").(string); v != "" {
		formData["tags"] = v
	}

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/software/")
	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42SoftwareCreate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42SoftwareCreate - Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42SoftwareCreate - unable to extract software id from response: %s", err)
	}

	d.SetId(strconv.Itoa(id))
	d.Set("software_id", id)

	if d.Get("custom_fields") != nil {
		fields := d.Get("custom_fields").(map[string]interface{})
		if len(fields) > 0 {
			bulkFields := []string{}
			for k, v := range fields {
				bulkFields = append(bulkFields, fmt.Sprintf("%v:%v", k, v))
			}

			log.Printf("[DEBUG] resourceDevice42SoftwareCreate - Pushing custom_fields %s", bulkFields)
			resp, err := client.R().
				SetFormData(map[string]string{
					"name":        name,
					"bulk_fields": strings.Join(bulkFields, ","),
				}).
				SetResult(apiResponse{}).
				Put("/1.0/custom_fields/software/")
			if err != nil {
				return err
			}

			r := resp.Result().(*apiResponse)
			if r.Code != 0 {
				return fmt.Errorf("resourceDevice42SoftwareCreate - api returned code %d", r.Code)
			}
		}
	}

	return resourceDevice42SoftwareRead(d, m)
}

func resourceDevice42SoftwareRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42SoftwareRead - starting read for software id %s", d.Id())

	var resp apiSoftwareListResponse
	_, err := client.R().
		SetResult(&resp).
		Get("/1.0/software/")
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("resourceDevice42SoftwareRead - invalid software id %q: %s", d.Id(), err)
	}

	for _, s := range resp.Software {
		if int(s.ID) == id {
			d.Set("name", s.Name)
			d.Set("software_id", s.ID)
			d.Set("vendor", s.Vendor)
			d.Set("description", s.Description)
			d.Set("notes", s.Notes)
			d.Set("software_type", s.SoftwareType)
			d.Set("deployment_type", s.DeploymentType)
			d.Set("category", s.Category)
			d.Set("aliases", s.Aliases)
			d.Set("licensing_model", s.LicensingModel)
			d.Set("track_licensed_count_by_keys", s.TrackLicensedCountByKeys)
			d.Set("tags", joinBuildingList(s.Tags))
			d.Set("custom_fields", flattenCustomFields(s.CustomFields))
			return nil
		}
	}

	log.Printf("[WARN] resourceDevice42SoftwareRead - no software found: %s", d.Id())
	d.SetId("")
	return nil
}

func resourceDevice42SoftwareUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	name := d.Get("name").(string)
	formData := map[string]string{}
	formData["name"] = name

	if v := d.Get("vendor").(string); v != "" {
		formData["vendor"] = v
	}
	if v := d.Get("description").(string); v != "" {
		formData["description"] = v
	}
	if v := d.Get("notes").(string); v != "" {
		formData["notes"] = v
	}
	if v := d.Get("software_type").(string); v != "" {
		formData["software_type"] = v
	}
	if v := d.Get("deployment_type").(string); v != "" {
		formData["deployment_type"] = v
	}
	if v := d.Get("category").(string); v != "" {
		formData["category"] = v
	}
	if v := d.Get("aliases").(string); v != "" {
		formData["aliases"] = v
	}
	if v := d.Get("licensing_model").(string); v != "" {
		formData["licensing_model"] = v
	}
	if v := d.Get("track_licensed_count_by_keys").(string); v != "" {
		formData["track_licensed_count_by_keys"] = v
	}
	if v := d.Get("tags").(string); v != "" {
		formData["tags"] = v
	}

	log.Printf("[DEBUG] resourceDevice42SoftwareUpdate - Pushing new informations : %#v", formData)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/software/")
	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42SoftwareUpdate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42SoftwareUpdate - Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42SoftwareUpdate - unable to extract software id from response: %s", err)
	}

	d.SetId(strconv.Itoa(id))
	d.Set("software_id", id)

	if d.HasChange("custom_fields") {
		updateList := setCustomFields(d)
		for k, v := range updateList {
			log.Printf("[DEBUG] resourceDevice42SoftwareUpdate - Pushing custom field %s : %v", k, v)
			resp, err := client.R().
				SetFormData(map[string]string{
					"name":  name,
					"key":   k,
					"value": v.(string),
				}).
				SetResult(apiResponse{}).
				Put("/1.0/custom_fields/software/")
			if err != nil {
				return err
			}

			r := resp.Result().(*apiResponse)
			log.Printf("[DEBUG] resourceDevice42SoftwareUpdate - Result: %#v", r)
		}
	}

	return resourceDevice42SoftwareRead(d, m)
}

func resourceDevice42SoftwareDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting software %s (ID: %s)", d.Get("name"), d.Id())

	url := fmt.Sprintf("/1.0/software/%s/", d.Id())

	resp, err := client.R().
		SetResult(apiDeleteResponse{}).
		Delete(url)
	if err != nil {
		return err
	}

	r := resp.Result().(*apiDeleteResponse)
	log.Printf("[DEBUG] resourceDevice42SoftwareDelete - Result: %#v", r)
	return nil
}
