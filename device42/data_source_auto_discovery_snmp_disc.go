package device42

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42AutoDiscoverySnmpDisc() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42AutoDiscoverySnmpDiscRead,
		Description: "Read a Device42 SNMP discovery autodiscovery job.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the SNMP discovery autodiscovery job.",
			},
			"job_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the SNMP discovery autodiscovery job.",
			},
			"server": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"exclude_servers": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snmp_retries": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"snmp_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"snmp_min_parallel": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"do_port_check": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"clear_customer": {
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
			"snmp_v3": {
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

func datasourceD42AutoDiscoverySnmpDiscRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	jobID := d.Get("job_id").(int)

	if name == "" && jobID == 0 {
		return fmt.Errorf("either 'name' or 'job_id' must be set to look up a Device42 SNMP discovery autodiscovery job")
	}

	var resp apiAutoDiscoverySnmpDiscListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/auto_discovery/snmp_disc/")
	if err != nil {
		log.Printf("[WARN] No snmp_disc autodiscovery job found")
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("datasourceD42AutoDiscoverySnmpDiscRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	var match *apiAutoDiscoverySnmpDisc
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
		log.Printf("[ERROR] No snmp_disc autodiscovery job found with name: %s job_id: %d", name, jobID)
		d.SetId("")
		return nil
	}

	d.SetId(strconv.FormatInt(match.JobId, 10))
	d.Set("job_id", match.JobId)
	d.Set("name", match.Name)
	d.Set("server", match.Server)
	d.Set("start_ip", match.StartIp)
	d.Set("end_ip", match.EndIp)
	d.Set("exclude_servers", match.ExcludeServers)
	d.Set("snmp_retries", match.SnmpRetries)
	d.Set("snmp_timeout", match.SnmpTimeout)
	d.Set("snmp_min_parallel", match.SnmpMinParallel)
	d.Set("do_port_check", fmt.Sprintf("%v", match.DoPortCheck))
	d.Set("clear_customer", fmt.Sprintf("%v", match.ClearCustomer))
	d.Set("schedule_time", match.ScheduleTime)
	d.Set("schedule_days", match.ScheduleDays)
	d.Set("snmp_v3", match.SnmpV3)
	d.Set("tags", joinBuildingList(match.Tags))

	return nil
}
