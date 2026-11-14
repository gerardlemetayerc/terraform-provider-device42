package device42

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42PduModel() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42PduModelRead,
		Description: "Read a Device42 PDU model.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the PDU model.",
			},
			"pdu_model_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the PDU model.",
			},
			"manufacturer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"depth": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"width": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"height": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func datasourceD42PduModelRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	modelID := d.Get("pdu_model_id").(int)

	if name == "" && modelID == 0 {
		return fmt.Errorf("either 'name' or 'pdu_model_id' must be set to look up a Device42 PDU model")
	}

	var resp apiPduModelListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/pdu_models/")
	if err != nil {
		log.Printf("[WARN] No pdu model found")
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("datasourceD42PduModelRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	var match *apiPduModel
	for i := range resp.PduModels {
		m := resp.PduModels[i]
		if modelID > 0 && int(m.PduModelId) == modelID {
			match = &resp.PduModels[i]
			break
		}
		if name != "" && m.Name == name {
			match = &resp.PduModels[i]
			break
		}
	}
	if match == nil {
		log.Printf("[ERROR] No pdu model found with name: %s id: %d", name, modelID)
		d.SetId("")
		return nil
	}

	d.SetId(strconv.FormatInt(match.PduModelId, 10))
	d.Set("pdu_model_id", match.PduModelId)
	d.Set("name", match.Name)
	d.Set("manufacturer", match.Manufacturer)
	d.Set("description", match.Description)
	d.Set("notes", match.Notes)
	d.Set("depth", match.Depth)
	d.Set("size", int64(match.Size))
	d.Set("width", int64(match.Width))
	d.Set("height", int64(match.Height))

	return nil
}
