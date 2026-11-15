package device42

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42AutoDiscoveryDns() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42AutoDiscoveryDnsRead,
		Description: "Read a Device42 DNS autodiscovery job.",
		Schema: map[string]*schema.Schema{
			"zonename": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zone name of the DNS autodiscovery job.",
			},
			"job_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the DNS autodiscovery job.",
			},
			"nameserver": {
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
		},
	}
}

func datasourceD42AutoDiscoveryDnsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	zonename := d.Get("zonename").(string)
	jobID := d.Get("job_id").(int)

	if zonename == "" && jobID == 0 {
		return fmt.Errorf("either 'zonename' or 'job_id' must be set to look up a Device42 DNS autodiscovery job")
	}

	var resp apiAutoDiscoveryDnsListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/auto_discovery/dns/")
	if err != nil {
		log.Printf("[WARN] No DNS autodiscovery job found")
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("datasourceD42AutoDiscoveryDnsRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	var match *apiAutoDiscoveryDns
	for i := range resp.Jobs {
		j := resp.Jobs[i]
		if jobID > 0 && int(j.JobId) == jobID {
			match = &resp.Jobs[i]
			break
		}
		if zonename != "" && j.ZoneName == zonename {
			match = &resp.Jobs[i]
			break
		}
	}
	if match == nil {
		log.Printf("[ERROR] No DNS autodiscovery job found with zonename: %s job_id: %d", zonename, jobID)
		d.SetId("")
		return nil
	}

	d.SetId(strconv.FormatInt(match.JobId, 10))
	d.Set("job_id", match.JobId)
	d.Set("zonename", match.ZoneName)
	d.Set("nameserver", match.Nameserver)
	d.Set("schedule_time", match.ScheduleTime)
	d.Set("schedule_days", match.ScheduleDays)

	return nil
}
