package device42

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiPduListResponse struct {
	Pdus []apiPdu `json:"pdus"`
}

func datasourceD42Pdu() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42PduRead,
		Description: "Read a Device42 PDU.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the PDU.",
			},
			"pdu_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the PDU.",
			},
			"pdu_serial_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pdu_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pdu_model_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway_pdu_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_fields": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func datasourceD42PduRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	pduID := d.Get("pdu_id").(int)

	if name == "" && pduID == 0 {
		return fmt.Errorf("either 'name' or 'pdu_id' must be set to look up a Device42 PDU")
	}

	var resp apiPduListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/pdus/")
	if err != nil {
		log.Printf("[WARN] No pdu found")
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("datasourceD42PduRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	var match *apiPdu
	for i := range resp.Pdus {
		p := resp.Pdus[i]
		if pduID > 0 && int(p.PduId) == pduID {
			match = &resp.Pdus[i]
			break
		}
		if name != "" && p.Name == name {
			match = &resp.Pdus[i]
			break
		}
	}
	if match == nil {
		log.Printf("[ERROR] No pdu found with name: %s id: %d", name, pduID)
		d.SetId("")
		return nil
	}

	d.SetId(strconv.FormatInt(match.PduId, 10))
	d.Set("pdu_id", match.PduId)
	d.Set("name", match.Name)
	d.Set("pdu_serial_no", match.PduSerialNo)
	if len(match.PduModel) > 0 && string(match.PduModel) != `""` {
		var model apiPduModel
		if err := json.Unmarshal(match.PduModel, &model); err != nil {
			return fmt.Errorf("datasourceD42PduRead - malformed pdu_model: %v", err)
		}
		d.Set("pdu_model", model.Name)
		d.Set("pdu_model_id", strconv.FormatInt(model.PduModelId, 10))
	}
	d.Set("notes", match.Notes)
	d.Set("gateway_pdu_id", match.GatewayPduId)
	d.Set("custom_fields", flattenCustomFields(match.CustomFields))

	return nil
}
