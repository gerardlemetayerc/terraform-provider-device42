package device42

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42AutoDiscoveryPingSweep() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42AutoDiscoveryPingSweepRead,
		Description: "Read a Device42 ping sweep autodiscovery job.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the ping sweep autodiscovery job.",
			},
			"job_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the ping sweep autodiscovery job.",
			},
			"networks": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schedule_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_new_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"debug_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"overwrite_subnet_categories": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schedule_days": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reverse_dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"strip_domain_suffix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceD42AutoDiscoveryPingSweepRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	jobID := d.Get("job_id").(int)

	if name == "" && jobID == 0 {
		return fmt.Errorf("either 'name' or 'job_id' must be set to look up a Device42 ping sweep autodiscovery job")
	}

	var resp apiAutoDiscoveryPingSweepListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/auto_discovery/pingsweep/")
	if err != nil {
		log.Printf("[WARN] No ping sweep autodiscovery job found")
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("datasourceD42AutoDiscoveryPingSweepRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	var match *apiAutoDiscoveryPingSweep
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
		log.Printf("[ERROR] No ping sweep autodiscovery job found with name: %s job_id: %d", name, jobID)
		d.SetId("")
		return nil
	}

	d.SetId(strconv.FormatInt(match.JobId, 10))
	d.Set("job_id", match.JobId)
	d.Set("name", match.Name)
	d.Set("networks", match.Networks)
	d.Set("schedule_time", match.ScheduleTime)
	d.Set("subnet_category", match.SubnetCategory)
	d.Set("create_new_subnet", match.CreateNewSubnet)
	d.Set("debug_level", match.DebugLevel)
	d.Set("overwrite_subnet_categories", match.OverwriteSubnetCategories)
	d.Set("schedule_days", match.ScheduleDays)
	d.Set("reverse_dns", match.ReverseDns)
	d.Set("strip_domain_suffix", match.StripDomainSuffix)
	d.Set("type", match.Type)

	return nil
}
