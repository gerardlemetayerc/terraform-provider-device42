package device42

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42AutoDiscoveryIpmi() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42AutoDiscoveryIpmiRead,
		Description: "Read a Device42 IPMI/Redfish autodiscovery job.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the IPMI/Redfish autodiscovery job.",
			},
			"job_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the IPMI/Redfish autodiscovery job.",
			},
			"server": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schedule_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_model_if_found": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"overwrite_object_categories": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_start": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_end": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"overwrite_device_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hostname_to_use": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"debug_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"discovery_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"run_as_operator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schedule_days": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"exclude_servers": {
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

func datasourceD42AutoDiscoveryIpmiRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	jobID := d.Get("job_id").(int)

	if name == "" && jobID == 0 {
		return fmt.Errorf("either 'name' or 'job_id' must be set to look up a Device42 IPMI/Redfish autodiscovery job")
	}

	var resp apiAutoDiscoveryIpmiListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/auto_discovery/ipmi/")
	if err != nil {
		log.Printf("[WARN] No IPMI/Redfish autodiscovery job found")
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("datasourceD42AutoDiscoveryIpmiRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	var match *apiAutoDiscoveryIpmi
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
		log.Printf("[ERROR] No IPMI/Redfish autodiscovery job found with name: %s job_id: %d", name, jobID)
		d.SetId("")
		return nil
	}

	d.SetId(strconv.FormatInt(match.JobId, 10))
	d.Set("job_id", match.JobId)
	d.Set("name", match.Name)
	d.Set("server", match.Server)
	d.Set("schedule_time", match.ScheduleTime)
	d.Set("object_category", match.ObjectCategory)
	d.Set("update_model_if_found", match.UpdateModelIfFound)
	d.Set("overwrite_object_categories", match.OverwriteObjectCategories)
	d.Set("ip_start", match.IpStart)
	d.Set("ip_end", match.IpEnd)
	d.Set("overwrite_device_name", match.OverwriteDeviceName)
	d.Set("hostname_to_use", match.HostnameToUse)
	d.Set("debug_level", match.DebugLevel)
	d.Set("discovery_type", match.DiscoveryType)
	d.Set("run_as_operator", match.RunAsOperator)
	d.Set("schedule_days", match.ScheduleDays)
	d.Set("exclude_servers", match.ExcludeServers)
	d.Set("tags", joinBuildingList(match.Tags))

	return nil
}
