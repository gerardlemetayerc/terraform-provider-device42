package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiAutoDiscoveryCloudAccount struct {
	JobId                      int64         `json:"job_id"`
	Name                       string        `json:"name"`
	Type                       string        `json:"type"`
	Server                     string        `json:"server"`
	IP                         string        `json:"ip"`
	Username                   string        `json:"username"`
	Tenant                     string        `json:"tenant"`
	Vendor                     string        `json:"vendor"`
	ObjectCategory             string        `json:"object_category"`
	VrfGroup                   string        `json:"vrf_group"`
	Notes                      string        `json:"notes"`
	ServiceLevel               string        `json:"service_level"`
	Regions                    string        `json:"regions"`
	AwsRoles                   string        `json:"aws_roles"`
	OracleRegions              string        `json:"oracle_regions"`
	GoogleRegions              string        `json:"google_regions"`
	AlibabaRegions             string        `json:"alibaba_regions"`
	OutscaleRegions            string        `json:"outscale_regions"`
	CollectDbInfo              string        `json:"collect_db_info"`
	AddSuffix                  string        `json:"add_suffix"`
	StripDomainName            string        `json:"strip_domain_name"`
	MatchNameOnlyForVirtuals   string        `json:"match_name_only_for_virtuals"`
	MatchNameOnlyForHypervisor string        `json:"match_name_only_for_hypervisor"`
	OverwriteObjectCategory    string        `json:"overwrite_object_category"`
	RemoveUnfoundInstances     string        `json:"remove_unfound_instances"`
	AddTagsAsCustomFields      string        `json:"add_tags_as_custom_fields"`
	ScheduleTime               string        `json:"schedule_time"`
	ScheduleDays               string        `json:"schedule_days"`
	Tags                       []interface{} `json:"tags"`
}

type apiAutoDiscoveryCloudAccountListResponse struct {
	Jobs []apiAutoDiscoveryCloudAccount `json:"jobs"`
}

func resourceD42AutoDiscoveryCloudAccount() *schema.Resource {
	return &schema.Resource{
		Description: "device42_auto_discovery_cloudaccount can be use to manage a cloud account autodiscovery job.",
		Create:      resourceDevice42AutoDiscoveryCloudAccountCreate,
		Read:        resourceDevice42AutoDiscoveryCloudAccountRead,
		Update:      resourceDevice42AutoDiscoveryCloudAccountUpdate,
		Delete:      resourceDevice42AutoDiscoveryCloudAccountDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"server": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"object_category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vrf_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"regions": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_roles": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oracle_regions": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"google_regions": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alibaba_regions": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"outscale_regions": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"collect_db_info": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"add_suffix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"strip_domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"match_name_only_for_virtuals": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"match_name_only_for_hypervisor": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"overwrite_object_category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remove_unfound_instances": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"add_tags_as_custom_fields": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule_days": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressCommaSeparatedOrderDiff,
			},
		},
	}
}

func resourceDevice42AutoDiscoveryCloudAccountCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := autoDiscoveryCloudAccountFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/auto_discovery/cloudaccount/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryCloudAccountCreate - API returned code %d: %v", r.Code, r.Msg)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryCloudAccountCreate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42AutoDiscoveryCloudAccountRead(d, m)
}

func resourceDevice42AutoDiscoveryCloudAccountRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42AutoDiscoveryCloudAccountRead - Starting reading using API for id %s", d.Id())

	var resp apiAutoDiscoveryCloudAccountListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/auto_discovery/cloudaccount/")
	if err != nil {
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryCloudAccountRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryCloudAccountRead - invalid id %q: %s", d.Id(), err)
	}

	for _, j := range resp.Jobs {
		if int(j.JobId) == id {
			d.Set("name", j.Name)
			d.Set("type", j.Type)
			d.Set("server", j.Server)
			d.Set("ip", j.IP)
			d.Set("username", j.Username)
			d.Set("tenant", j.Tenant)
			d.Set("vendor", j.Vendor)
			d.Set("object_category", j.ObjectCategory)
			d.Set("vrf_group", j.VrfGroup)
			d.Set("notes", j.Notes)
			d.Set("service_level", j.ServiceLevel)
			d.Set("regions", j.Regions)
			d.Set("aws_roles", j.AwsRoles)
			d.Set("oracle_regions", j.OracleRegions)
			d.Set("google_regions", j.GoogleRegions)
			d.Set("alibaba_regions", j.AlibabaRegions)
			d.Set("outscale_regions", j.OutscaleRegions)
			d.Set("collect_db_info", j.CollectDbInfo)
			d.Set("add_suffix", j.AddSuffix)
			d.Set("strip_domain_name", j.StripDomainName)
			d.Set("match_name_only_for_virtuals", j.MatchNameOnlyForVirtuals)
			d.Set("match_name_only_for_hypervisor", j.MatchNameOnlyForHypervisor)
			d.Set("overwrite_object_category", j.OverwriteObjectCategory)
			d.Set("remove_unfound_instances", j.RemoveUnfoundInstances)
			d.Set("add_tags_as_custom_fields", j.AddTagsAsCustomFields)
			d.Set("schedule_time", j.ScheduleTime)
			d.Set("schedule_days", j.ScheduleDays)
			d.Set("tags", joinBuildingList(j.Tags))
			return nil
		}
	}

	log.Printf("[WARN] resourceDevice42AutoDiscoveryCloudAccountRead - no cloud account job found: %s", d.Id())
	d.SetId("")
	return nil
}

func resourceDevice42AutoDiscoveryCloudAccountUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	formData := autoDiscoveryCloudAccountFormData(d)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/auto_discovery/cloudaccount/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryCloudAccountUpdate - API returned code %d: %v", r.Code, r.Msg)
	}
	log.Printf("[DEBUG] Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42AutoDiscoveryCloudAccountUpdate - malformed API response: %v", err)
	}
	d.SetId(strconv.Itoa(id))
	return resourceDevice42AutoDiscoveryCloudAccountRead(d, m)
}

func resourceDevice42AutoDiscoveryCloudAccountDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting cloud account autodiscovery job id %s (name: %s)", d.Id(), d.Get("name"))

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(fmt.Sprintf("/1.0/auto_discovery/cloudaccount/%s/", d.Id()))

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42AutoDiscoveryCloudAccountDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}

// autoDiscoveryCloudAccountFormData builds the create/update payload. The API
// upserts by name: re-posting the same name updates the existing job, so the
// form data is identical for create and update.
func autoDiscoveryCloudAccountFormData(d *schema.ResourceData) map[string]string {
	formData := map[string]string{
		"name": d.Get("name").(string),
		"type": d.Get("type").(string),
	}

	if v, ok := d.GetOk("server"); ok {
		formData["server"] = v.(string)
	}
	if v, ok := d.GetOk("ip"); ok {
		formData["ip"] = v.(string)
	}
	if v, ok := d.GetOk("username"); ok {
		formData["username"] = v.(string)
	}
	if v, ok := d.GetOk("tenant"); ok {
		formData["tenant"] = v.(string)
	}
	if v, ok := d.GetOk("vendor"); ok {
		formData["vendor"] = v.(string)
	}
	if v, ok := d.GetOk("object_category"); ok {
		formData["object_category"] = v.(string)
	}
	if v, ok := d.GetOk("vrf_group"); ok {
		formData["vrf_group"] = v.(string)
	}
	if v, ok := d.GetOk("notes"); ok {
		formData["notes"] = v.(string)
	}
	if v, ok := d.GetOk("service_level"); ok {
		formData["service_level"] = v.(string)
	}
	if v, ok := d.GetOk("regions"); ok {
		formData["regions"] = v.(string)
	}
	if v, ok := d.GetOk("aws_roles"); ok {
		formData["aws_roles"] = v.(string)
	}
	if v, ok := d.GetOk("oracle_regions"); ok {
		formData["oracle_regions"] = v.(string)
	}
	if v, ok := d.GetOk("google_regions"); ok {
		formData["google_regions"] = v.(string)
	}
	if v, ok := d.GetOk("alibaba_regions"); ok {
		formData["alibaba_regions"] = v.(string)
	}
	if v, ok := d.GetOk("outscale_regions"); ok {
		formData["outscale_regions"] = v.(string)
	}
	if v, ok := d.GetOk("collect_db_info"); ok {
		formData["collect_db_info"] = v.(string)
	}
	if v, ok := d.GetOk("add_suffix"); ok {
		formData["add_suffix"] = v.(string)
	}
	if v, ok := d.GetOk("strip_domain_name"); ok {
		formData["strip_domain_name"] = v.(string)
	}
	if v, ok := d.GetOk("match_name_only_for_virtuals"); ok {
		formData["match_name_only_for_virtuals"] = v.(string)
	}
	if v, ok := d.GetOk("match_name_only_for_hypervisor"); ok {
		formData["match_name_only_for_hypervisor"] = v.(string)
	}
	if v, ok := d.GetOk("overwrite_object_category"); ok {
		formData["overwrite_object_category"] = v.(string)
	}
	if v, ok := d.GetOk("remove_unfound_instances"); ok {
		formData["remove_unfound_instances"] = v.(string)
	}
	if v, ok := d.GetOk("add_tags_as_custom_fields"); ok {
		formData["add_tags_as_custom_fields"] = v.(string)
	}
	if v, ok := d.GetOk("schedule_time"); ok {
		formData["schedule_time"] = v.(string)
	}
	if v, ok := d.GetOk("schedule_days"); ok {
		formData["schedule_days"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		formData["tags"] = v.(string)
	}

	return formData
}
