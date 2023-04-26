package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)


type apiSubnetReadResponse struct {
	ID                      int64         `json:"id"`
	Network                 string 	      `json:"network"`	
	Mask_Bits               int64         `json:"mask_bits"`
	vrfGroup                string        `json:"vrf_group"`
	Name                    string        `json:"name"`
}



func resourceD42Subnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceDevice42SubnetCreate,
		Read:   resourceDevice42SubnetRead,
		Update: resourceDevice42SubnetUpdate,
		Delete: resourceDevice42SubnetDelete,

		Schema: map[string]*schema.Schema{
			"network": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Network of the subnet. Required for creation, cannot be modified after subnet creation.",
			},
			"mask_bits": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Mask bits of the subnet. Required for creation, can be modified after subnet creation.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The hostname of the device.",
			},
			"vrf_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet VRF Group",
			},
			"parent_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Parent subnet id",
			},
			"customer": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Any custom fields that will be used in device42.",
			},
			"description": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"gateway": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"service_level": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"vlan": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceDevice42SubnetCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	name := d.Get("name").(string)
	network := d.Get("network").(string)
	maskBits := d.Get("mask_bits").(string)
	vrfGroup := d.Get("vrf_group").(string)

	resp, err := client.R().
		SetFormData(map[string]string{
			"name":      name,
			"network":   network,
			"mask_bits": maskBits,
			"vrfGroup" : vrfGroup,
		}).
		SetResult(apiResponse{}).
		Post("/1.0/subnets/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("API returned code %d", r.Code)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id := int(r.Msg[1].(float64))

	// Set ID after subnet creation
	d.SetId(strconv.Itoa(id))

	return resourceDevice42DeviceRead(d, m)
}

func resourceDevice42SubnetRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	resp, err := client.R().
		SetResult(apiSubnetReadResponse{}).
		Get(fmt.Sprintf("/1.0/subnets/%s/", d.Id()))

	if err != nil {
		log.Printf("[WARN] No subnet found: %s", d.Id())
		d.SetId("")
		return err
	}

	r := resp.Result().(*apiSubnetReadResponse)

	d.Set("name", r.Name)
	d.Set("network", r.Network)
	d.Set("mask_bits", r.Mask_Bits)
	d.Set("vrf_group", r.vrfGroup)

	return nil
}

func resourceDevice42SubnetUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	if(d.HasChange("name") || d.HasChange("mask_bits") || d.HasChange("vrf_group")){
		name := d.Get("name").(string)
		maskBits := d.Get("mask_bits").(string)
		vrfGroup := d.Get("vrf_group").(string)
		url := fmt.Sprintf("/1.0/subnets/%s/", d.Id())


		resp, err := client.R().
			SetFormData(map[string]string{
				"name":        name,
				"mask_bits":   maskBits,
				"vrf_group":   vrfGroup,
			}).
			SetResult(apiResponse{}).
			Put(url)

		if err != nil {
			return err
		}
		r := resp.Result().(*apiResponse)
		log.Printf("[DEBUG] Result: %#v", r)
	}


	return resourceDevice42SubnetRead(d, m)
}

func resourceDevice42SubnetDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting subnet %s (UUID: %s)", d.Get("name"), d.Id())

	url := fmt.Sprintf("/1.0/subnets/%s/", d.Id())

	resp, err := client.R().
		SetResult(apiResponse{}).
		Delete(url)

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}
