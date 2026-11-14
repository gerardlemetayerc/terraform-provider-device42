package device42

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// apiMacListItem mirrors the shape returned by GET /1.0/macs/. The "port"
// field can be null, a plain string, or a nested object, so it is kept as
// json.RawMessage and rendered only when it is a plain string.
type apiMacListItem struct {
	MacaddressId int64           `json:"macaddress_id"`
	Macaddress   string          `json:"macaddress"`
	VlanId       int64           `json:"vlan_id"`
	Port         json.RawMessage `json:"port"`
}

type apiMacListResponse struct {
	Macaddresses []apiMacListItem `json:"macaddresses"`
}

func datasourceD42Mac() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42MacRead,
		Description: "Read a Device42 MAC address.",
		Schema: map[string]*schema.Schema{
			"macaddress": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "MAC address value.",
			},
			"macaddress_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the MAC address.",
			},
			"vlan_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceD42MacRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	macaddress := d.Get("macaddress").(string)
	macID := d.Get("macaddress_id").(int)

	if macaddress == "" && macID == 0 {
		return fmt.Errorf("either 'macaddress' or 'macaddress_id' must be set to look up a Device42 MAC address")
	}

	var resp apiMacListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/macs/")
	if err != nil {
		log.Printf("[WARN] No mac address found")
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("datasourceD42MacRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	var match *apiMacListItem
	for i := range resp.Macaddresses {
		m := resp.Macaddresses[i]
		if macID > 0 && int(m.MacaddressId) == macID {
			match = &resp.Macaddresses[i]
			break
		}
		if macaddress != "" && strings.EqualFold(m.Macaddress, macaddress) {
			match = &resp.Macaddresses[i]
			break
		}
	}
	if match == nil {
		log.Printf("[ERROR] No mac address found with macaddress: %s id: %d", macaddress, macID)
		d.SetId("")
		return nil
	}

	d.SetId(strconv.FormatInt(match.MacaddressId, 10))
	d.Set("macaddress_id", match.MacaddressId)
	d.Set("macaddress", match.Macaddress)
	d.Set("vlan_id", strconv.FormatInt(match.VlanId, 10))
	d.Set("port", renderMacPort(match.Port))

	return nil
}

// renderMacPort renders the "port" field of a MAC address list item. It can be
// null, a plain string, or a nested object; only a plain string is surfaced.
func renderMacPort(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" || string(raw) == `""` {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	return ""
}
