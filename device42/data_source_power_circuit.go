package device42

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiPowerCircuitListResponse struct {
	PowerCircuits []apiPowerCircuitResponse `json:"power_circuits"`
}

func datasourceD42PowerCircuit() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42PowerCircuitRead,
		Description: "Read a Device42 power circuit.",
		Schema: map[string]*schema.Schema{
			"power_circuit_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the power circuit.",
			},
			"breakerpanel_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifier of the breaker panel.",
			},
			"number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Number of the circuit within the breaker panel.",
			},
			"breakerpanel": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"amps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"voltage": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"connection_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceD42PowerCircuitRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	circuitID := d.Get("power_circuit_id").(int)
	breakerpanelID := d.Get("breakerpanel_id").(string)
	number := d.Get("number").(string)

	if circuitID == 0 && (breakerpanelID == "" || number == "") {
		return fmt.Errorf("either 'power_circuit_id', or both 'breakerpanel_id' and 'number', must be set to look up a Device42 power circuit")
	}

	var resp apiPowerCircuitListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/powercircuits/")
	if err != nil {
		log.Printf("[WARN] No power circuit found")
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("datasourceD42PowerCircuitRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	var match *apiPowerCircuitResponse
	for i := range resp.PowerCircuits {
		c := resp.PowerCircuits[i]
		if circuitID > 0 && int(c.PowerCircuitId) == circuitID {
			match = &resp.PowerCircuits[i]
			break
		}
		if breakerpanelID != "" && number != "" &&
			strconv.FormatInt(c.BreakerpanelId, 10) == breakerpanelID && c.Number == number {
			match = &resp.PowerCircuits[i]
			break
		}
	}
	if match == nil {
		log.Printf("[ERROR] No power circuit found with id: %d breakerpanel_id: %s number: %s", circuitID, breakerpanelID, number)
		d.SetId("")
		return nil
	}

	d.SetId(strconv.FormatInt(match.PowerCircuitId, 10))
	d.Set("power_circuit_id", match.PowerCircuitId)
	d.Set("breakerpanel_id", strconv.FormatInt(match.BreakerpanelId, 10))
	d.Set("number", match.Number)
	d.Set("breakerpanel", match.Breakerpanel)
	d.Set("label", match.Label)
	d.Set("amps", int64(match.Amps))
	d.Set("voltage", int64(match.Voltage))
	d.Set("connection_name", match.Connection)

	return nil
}
