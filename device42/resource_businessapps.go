package device42

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiBusinessAppsReadResponse struct {
	AppType              string        `json:"app_type"`
	AppTypeId            int64         `json:"app_type_id"`
	BusinessAppId        int64         `json:"businessapp_id"`
	BusinessAppOwner     string        `json:"business_app_owner"`
	BusinessAppOwerId    int64         `json:"business_app_owner_id"`
	Created              string        `json:"created"`
	CriticalityId        int64         `json:"criticality_id"`
	CustOwner            string        `json:"cust_owner"`
	CustOwnerId          int64         `json:"cust_owner_id"`
	Description          string        `json:"description"`
	IsContainsPII        bool          `json:"is_contains_pii"`
	IsInternetAccessible bool          `json:"is_internet_accessible"`
	LastChanged          string        `json:"last_changed"`
	MigrationGroup       string        `json:"migration_group"`
	MigrationGroupId     int64         `json:"migration_group_id"`
	Name                 string        `json:"name"`
	Notes                string        `json:"notes"`
	ServiceLevel         string        `json:"service_level"`
	ServiceLevelId       int64         `json:"service_level_id"`
	TechnicalAppOwnerId  int64         `json:"technical_app_owner_id"`
	TechnicalAppOwner    string        `json:"technical_app_owner"`
	VendorId             int64         `json:"vendor_id"`
	CustomFields         []customField `json:"custom_fields"`
}

func resourceD42BusinessApps() *schema.Resource {
	return &schema.Resource{
		Description: "device42_businessapp can be use to manage Business Applications",
		Create:      resourceDevice42BusinessAppsCreate,
		Update:      resourceDevice42BusinessAppsUpdate,
		Delete:      resourceDevice42BusinessAppsDelete,
		Read:        resourceDevice42BusinessAppsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Name of the Business Application. REQUIRED to create a new application.",
			},
			"app_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Business App type.",
			},
			"app_type_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Business App type.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Business Application.",
			},
			"business_app_owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Business application owner name.",
			},
			"business_app_owner_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Business application owner ID.",
			},
			"technical_app_owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Technical application owner name.",
			},
			"technical_app_owner_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Technical application owner ID.",
			},
			"cust_owner_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Customer owner ID.",
			},
			"service_level_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "D42 ID of service level name (do not use with service_level).",
			},
			"service_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Service Level name (do not use with service_level_id)",
			},
			"migration_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Business application Migration Group name (do not use with migration_group_id)",
			},
			"migration_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "D42 ID of business application Migration Group (do not use with migration_group)",
			},
			"vendor": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Business Application vendor",
			},
			"custom_fields": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Any custom fields that will be used in device42.",
			},
			"is_contains_pii": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the Business Application contains PII.",
			},
			"is_internet_accessible": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the Business Application is accessible from the internet.",
			},
			"last_changed": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last change date of the Business Application.",
			},
			"notes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Any additional notes for the Business Application.",
			},
		},
	}
}

func resourceDevice42BusinessAppsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	resp, err := client.R().
		SetResult(apiBusinessAppsReadResponse{}).
		Get(fmt.Sprintf("/1.0/businessapps/%s/", d.Id()))

	if err != nil {
		log.Printf("[WARN] No Business App found: %s", d.Id())
		d.SetId("")
		return err
	}
	if resp.StatusCode() == 404 {
		log.Printf("[WARN] No Business App found: %s", d.Id())
		d.SetId("")
		return nil
	}
	if resp.StatusCode() >= 400 {
		return fmt.Errorf("resourceDevice42BusinessAppsRead - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiBusinessAppsReadResponse)
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("business_app_owner_id", strconv.FormatInt(r.BusinessAppOwerId, 10))
	d.Set("business_app_owner", r.BusinessAppOwner)
	d.Set("app_type", r.AppType)
	d.Set("app_type_id", strconv.FormatInt(r.AppTypeId, 10))
	d.Set("cust_owner_id", strconv.FormatInt(r.CustOwnerId, 10))
	d.Set("is_contains_pii", r.IsContainsPII)
	d.Set("is_internet_accessible", r.IsInternetAccessible)
	d.Set("last_changed", r.LastChanged)
	d.Set("migration_group", r.MigrationGroup)
	d.Set("migration_group_id", strconv.FormatInt(r.MigrationGroupId, 10))
	d.Set("notes", r.Notes)
	d.Set("service_level", r.ServiceLevel)
	d.Set("service_level_id", strconv.FormatInt(r.ServiceLevelId, 10))
	d.Set("technical_app_owner", r.TechnicalAppOwner)
	d.Set("technical_app_owner_id", strconv.FormatInt(r.TechnicalAppOwnerId, 10))
	d.Set("vendor", strconv.FormatInt(r.VendorId, 10))
	fields := flattenCustomFields(r.CustomFields)
	d.Set("custom_fields", fields)
	return nil
}

func resourceDevice42BusinessAppsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	business_app_owner_id := d.Get("business_app_owner_id").(string)
	technical_app_owner_id := d.Get("technical_app_owner_id").(string)
	cust_owner_id := d.Get("cust_owner_id").(string)
	service_level_id := d.Get("service_level_id").(string)
	service_level := d.Get("service_level").(string)
	migration_group := d.Get("migration_group").(string)
	migration_group_id := d.Get("migration_group_id").(string)

	mapData := map[string]string{
		"name": name,
	}

	if description != "" {
		mapData["description"] = description
	}

	if business_app_owner_id != "" {
		mapData["business_app_owner_id"] = business_app_owner_id
	}

	if technical_app_owner_id != "" {
		mapData["technical_app_owner_id"] = technical_app_owner_id
	}
	if cust_owner_id != "" {
		mapData["cust_owner_id"] = cust_owner_id
	}
	if service_level_id != "" {
		mapData["service_level_id"] = service_level_id
	}
	if service_level != "" {
		mapData["service_level"] = service_level
	}
	if migration_group != "" {
		mapData["migration_group"] = migration_group
	}
	if migration_group_id != "" {
		mapData["migration_group_id"] = migration_group_id
	}

	resp, err := client.R().
		SetFormData(mapData).
		SetResult(apiResponse{}).
		Post("/1.0/businessapps/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("API returned code %d", r.Code)
	}

	// The demo API upserts by name: a create for an existing name returns
	// {code:0, msg:[]} rather than the usual [msg, id, ...]. Recover the ID
	// by name lookup so create stays idempotent across re-runs.
	if len(r.Msg) < 1 {
		id, err := d42BusinessAppIdByName(client, name)
		if err != nil {
			return err
		}
		d.SetId(strconv.FormatInt(id, 10))
		return resourceDevice42BusinessAppsRead(d, m)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id := int(r.Msg[1].(float64))

	// Set ID after Business App creation
	d.SetId(strconv.Itoa(id))

	return resourceDevice42BusinessAppsRead(d, m)
}

func d42BusinessAppIdByName(client *resty.Client, name string) (int64, error) {
	resp, err := client.R().
		SetResult(datasourceD42BusinessAppApi{}).
		Get(fmt.Sprintf("/1.0/businessapps/?name=%s", url.QueryEscape(name)))
	if err != nil {
		return 0, err
	}
	if resp.StatusCode() >= 400 {
		return 0, fmt.Errorf("lookup businessapp by name - HTTP %d: %s", resp.StatusCode(), resp.String())
	}
	r := resp.Result().(*datasourceD42BusinessAppApi)
	if len(r.Businessapps) == 0 {
		return 0, fmt.Errorf("businessapp %q not found after create (empty API response)", name)
	}
	return r.Businessapps[0].BusinessAppId, nil
}

func resourceDevice42BusinessAppsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting Business Apps %s (UUID: %s)", d.Get("name"), d.Id())

	url := fmt.Sprintf("/1.0/businessapps/%s/", d.Id())

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(url)

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42BusinessAppsDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

func resourceDevice42BusinessAppsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := map[string]string{
		"id": d.Id(),
	}

	if d.HasChange("name") {
		formData["name"] = d.Get("name").(string)
	}
	if d.HasChange("description") {
		formData["description"] = d.Get("description").(string)
	}
	if d.HasChange("business_app_owner_id") {
		formData["business_app_owner_id"] = d.Get("business_app_owner_id").(string)
	}
	if d.HasChange("technical_app_owner_id") {
		formData["technical_app_owner_id"] = d.Get("technical_app_owner_id").(string)
	}
	if d.HasChange("cust_owner_id") {
		formData["cust_owner_id"] = d.Get("cust_owner_id").(string)
	}
	if d.HasChange("cust_owner_id") {
		formData["cust_owner_id"] = d.Get("cust_owner_id").(string)
	}
	if d.HasChange("service_level_id") {
		formData["service_level_id"] = d.Get("service_level_id").(string)
	}
	if d.HasChange("service_level") {
		formData["service_level"] = d.Get("service_level").(string)
	}
	if d.HasChange("migration_group") {
		formData["migration_group"] = d.Get("migration_group").(string)
	}
	if d.HasChange("migration_group_id") {
		formData["migration_group_id"] = d.Get("migration_group_id").(string)
	}

	if len(formData) > 1 {
		resp, err := client.R().
			SetFormData(formData).
			SetResult(apiResponse{}).
			Post("/1.0/businessapps/")

		if err != nil {
			return err
		}
		r := resp.Result().(*apiResponse)
		log.Printf("[DEBUG] Result: %#v", r)
	}

	return resourceDevice42BusinessAppsRead(d, m)
}
