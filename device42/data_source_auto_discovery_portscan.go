package device42

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42AutoDiscoveryPortScan() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42AutoDiscoveryPortScanRead,
		Description: "Read a Device42 tcp port scan autodiscovery job.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the tcp port scan autodiscovery job.",
			},
			"job_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the tcp port scan autodiscovery job.",
			},
			"start_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"capture_hosts_file": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"discover_cluster": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"overwrite_device_categories": {
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

func datasourceD42AutoDiscoveryPortScanRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	jobID := d.Get("job_id").(int)

	if name == "" && jobID == 0 {
		return fmt.Errorf("either 'name' or 'job_id' must be set to look up a Device42 tcp port scan autodiscovery job")
	}

	var resp apiAutoDiscoveryPortScanListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/auto_discovery/portscan/")
	if err != nil {
		log.Printf("[WARN] No tcp port scan autodiscovery job found")
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("datasourceD42AutoDiscoveryPortScanRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	var match *apiAutoDiscoveryPortScan
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
		log.Printf("[ERROR] No tcp port scan autodiscovery job found with name: %s job_id: %d", name, jobID)
		d.SetId("")
		return nil
	}

	d.SetId(strconv.FormatInt(match.JobId, 10))
	d.Set("job_id", match.JobId)
	d.Set("name", match.Name)
	d.Set("server", match.Server)
	d.Set("start_ip", match.StartIp)
	d.Set("end_ip", match.EndIp)
	d.Set("capture_hosts_file", match.CaptureHostsFile)
	d.Set("discover_cluster", match.DiscoverCluster)
	d.Set("overwrite_device_categories", match.OverwriteDeviceCategories)
	d.Set("schedule_time", match.ScheduleTime)
	d.Set("schedule_days", match.ScheduleDays)
	d.Set("tags", joinBuildingList(match.Tags))

	return nil
}
