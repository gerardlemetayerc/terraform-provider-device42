package device42

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42IpNat() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42IpNatRead,
		Description: "Read a Device42 IP NAT rule.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the IP NAT rule.",
			},
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the IP NAT rule.",
			},
			"ip_address_from": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address_to": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"two_way_relation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_port_start": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_port_end": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_port_start": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_port_end": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address_from_end": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address_to_end": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceD42IpNatRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	id := d.Get("id").(int)

	if name == "" && id == 0 {
		return fmt.Errorf("either 'name' or 'id' must be set to look up a Device42 IP NAT rule")
	}

	var resp apiIpNatListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/ipnat/")
	if err != nil {
		log.Printf("[WARN] No ip nat rule found")
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("datasourceD42IpNatRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	var match *apiIpNat
	for i := range resp.Ipnats {
		n := resp.Ipnats[i]
		if id > 0 && int(n.Id) == id {
			match = &resp.Ipnats[i]
			break
		}
		if name != "" && n.Name == name {
			match = &resp.Ipnats[i]
			break
		}
	}
	if match == nil {
		log.Printf("[ERROR] No ip nat rule found with name: %s id: %d", name, id)
		d.SetId("")
		return nil
	}

	d.SetId(strconv.FormatInt(match.Id, 10))
	d.Set("id", match.Id)
	d.Set("name", match.Name)
	d.Set("ip_address_from", match.IpAddressFrom)
	d.Set("ip_address_to", match.IpAddressTo)
	d.Set("notes", match.Notes)
	d.Set("protocol", match.Protocol)
	d.Set("two_way_relation", match.TwoWayRelation)
	d.Set("source_port_start", match.SourcePortStart)
	d.Set("source_port_end", match.SourcePortEnd)
	d.Set("target_port_start", match.TargetPortStart)
	d.Set("target_port_end", match.TargetPortEnd)
	d.Set("ip_address_from_end", match.IpAddressFromEnd)
	d.Set("ip_address_to_end", match.IpAddressToEnd)

	return nil
}
