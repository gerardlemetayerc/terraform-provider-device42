package device42

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type apiVlanReadReponse struct {
	Number      int64  `json:"number"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SwitchIds   string `json:"switch_ids"`
	Switches    string `json:"switches"`
	Tags        string `json:"tags"`
}

func resourceD42Vlans() *schema.Resource {
	return &schema.Resource{
		Description: "device42_businessapp can be use to manage Business Applications",
		Create:      resourceDevice42DeviceCreate,
		Read:        resourceDevice42DeviceRead,
		Update:      resourceDevice42DeviceUpdate,
		Delete:      resourceDevice42DeviceDelete,

		Schema: map[string]*schema.Schema{
			"number": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: "VLAN ID Number.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If left blank, name will be created as VLANxxxx, e.g. VLAN# 342 will be named VLAN0342",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"switches_ids": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comma separated values for switch id's",
			},
			"switches": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comma separated values for switch names",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Add or update tags to a VLAN",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Any additional notes",
			},
		},
	}
}
