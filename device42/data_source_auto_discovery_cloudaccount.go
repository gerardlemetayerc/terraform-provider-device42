package device42

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42AutoDiscoveryCloudAccount() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42AutoDiscoveryCloudAccountRead,
		Description: "Read a Device42 cloud account autodiscovery job.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the cloud account autodiscovery job.",
			},
			"job_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the cloud account autodiscovery job.",
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vrf_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"regions": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_roles": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"oracle_regions": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"google_regions": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alibaba_regions": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"outscale_regions": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"collect_db_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"add_suffix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"strip_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"match_name_only_for_virtuals": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"match_name_only_for_hypervisor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"overwrite_object_category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remove_unfound_instances": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"add_tags_as_custom_fields": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schedule_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schedule_days": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceD42AutoDiscoveryCloudAccountRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	jobID := d.Get("job_id").(int)

	if name == "" && jobID == 0 {
		return fmt.Errorf("either 'name' or 'job_id' must be set to look up a Device42 cloud account autodiscovery job")
	}

	var resp apiAutoDiscoveryCloudAccountListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/auto_discovery/cloudaccount/")
	if err != nil {
		log.Printf("[WARN] No cloud account autodiscovery job found")
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("datasourceD42AutoDiscoveryCloudAccountRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	var match *apiAutoDiscoveryCloudAccount
	for i := range resp.Jobs {
		j := resp.Jobs[i]
		if jobID > 0 && int(j.JobId) == jobID {
			match = &resp.Jobs[i]
			break
		}
		if name != "" && j.Name == name {
			match = &resp.Jobs[i]
			break
		}
	}
	if match == nil {
		log.Printf("[ERROR] No cloud account autodiscovery job found with name: %s job_id: %d", name, jobID)
		d.SetId("")
		return nil
	}

	d.SetId(strconv.FormatInt(match.JobId, 10))
	d.Set("job_id", match.JobId)
	d.Set("name", match.Name)
	d.Set("type", match.Type)
	d.Set("server", match.Server)
	d.Set("ip", match.IP)
	d.Set("username", match.Username)
	d.Set("tenant", match.Tenant)
	d.Set("vendor", match.Vendor)
	d.Set("object_category", match.ObjectCategory)
	d.Set("vrf_group", match.VrfGroup)
	d.Set("notes", match.Notes)
	d.Set("service_level", match.ServiceLevel)
	d.Set("regions", match.Regions)
	d.Set("aws_roles", match.AwsRoles)
	d.Set("oracle_regions", match.OracleRegions)
	d.Set("google_regions", match.GoogleRegions)
	d.Set("alibaba_regions", match.AlibabaRegions)
	d.Set("outscale_regions", match.OutscaleRegions)
	d.Set("collect_db_info", match.CollectDbInfo)
	d.Set("add_suffix", match.AddSuffix)
	d.Set("strip_domain_name", match.StripDomainName)
	d.Set("match_name_only_for_virtuals", match.MatchNameOnlyForVirtuals)
	d.Set("match_name_only_for_hypervisor", match.MatchNameOnlyForHypervisor)
	d.Set("overwrite_object_category", match.OverwriteObjectCategory)
	d.Set("remove_unfound_instances", match.RemoveUnfoundInstances)
	d.Set("add_tags_as_custom_fields", match.AddTagsAsCustomFields)
	d.Set("schedule_time", match.ScheduleTime)
	d.Set("schedule_days", match.ScheduleDays)
	d.Set("tags", joinBuildingList(match.Tags))

	return nil
}
