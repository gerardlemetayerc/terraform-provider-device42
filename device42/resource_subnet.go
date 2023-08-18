package device42

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiSubnetReadResponse struct {
	Subnet_id      int64         `json:"subnet_id"`
	Allocated      string        `json:"allocated"`
	Description    string        `json:"description"`
	Gateway        string        `json:"gateway"`
	MaskBits       int64         `json:"mask_bits"`
	Name           string        `json:"name"`
	Network        string        `json:"network"`
	RangeBegin     string        `json:"range_begin"`
	RangeEnd       string        `json:"range_end"`
	VrfGroupName   string        `json:"vrf_group_name"`
	VrfGroupId     int32         `json:"vrf_group_id"`
	ParentSubnetId int32         `json:"parent_subnet_id"`
	Customer       string        `json:"customer"`
	CustomFields   []customField `json:"custom_fields"`
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
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Mask bits of the subnet. Required for creation, can be modified after subnet creation.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet name.",
			},
			"vrf_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet VRF Group",
			},
			"parent_subnet_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Parent subnet id",
			},
			"customer": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Customer attached to this network.",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vlan": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vrf_group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"custom_fields": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				Description:      "Any custom fields that will be used in device42.",
				DiffSuppressFunc: suppressCustomFieldsDiffs,
			},
		},
	}
}

func resourceDevice42SubnetCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	resourceDevice42SubnetCreateForm := map[string]string{}
	if d.Get("name").(string) != "" {
		resourceDevice42SubnetCreateForm["name"] = d.Get("name").(string)
	}
	if d.Get("network").(string) != "" {
		resourceDevice42SubnetCreateForm["network"] = d.Get("network").(string)
	}

	if d.Get("mask_bits").(int) > 0 {
		resourceDevice42SubnetCreateForm["mask_bits"] = strconv.Itoa(d.Get("mask_bits").(int))
	}

	if d.Get("vrf_group").(string) != "" {
		resourceDevice42SubnetCreateForm["vrf_group"] = d.Get("vrf_group").(string)
	}
	if d.Get("gateway").(string) != "" {
		resourceDevice42SubnetCreateForm["gateway"] = d.Get("gateway").(string)
	}
	if d.Get("service_level").(string) != "" {
		resourceDevice42SubnetCreateForm["service_level"] = d.Get("service_level").(string)
	}
	if d.Get("description").(string) != "" {
		resourceDevice42SubnetCreateForm["description"] = d.Get("description").(string)
	}
	if d.Get("customer").(string) != "" {
		resourceDevice42SubnetCreateForm["customer"] = d.Get("customer").(string)
	}
	log.Printf("[DEBUG] vrf_group: %s", d.Get("vrf_group").(string))
	resp, err := client.R().
		SetFormData(resourceDevice42SubnetCreateForm).
		SetResult(apiResponse{}).
		Post("/1.0/subnets/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42SubnetCreate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42SubnetCreate - Result: %#v", r)
	if len(r.Msg) > 0 {
		id := int(r.Msg[1].(float64))
		// Set ID after subnet creation
		d.SetId(strconv.Itoa(id))
		if d.Get("custom_fields") != nil {
			fields := d.Get("custom_fields").(map[string]interface{})
			bulkFields := []string{}

			for k, v := range fields {
				bulkFields = append(bulkFields, fmt.Sprintf("%v:%v", k, v))
			}

			resp, err := client.R().
				SetFormData(map[string]string{
					"network":     resourceDevice42SubnetCreateForm["network"],
					"mask_bits":   resourceDevice42SubnetCreateForm["mask_bits"],
					"vrf_group":   resourceDevice42SubnetCreateForm["vrf_group"],
					"bulk_fields": strings.Join(bulkFields, ","),
				}).
				SetResult(apiResponse{}).
				Put("/1.0/custom_fields/subnet/")

			if err != nil {
				return err
			}

			r := resp.Result().(*apiResponse)

			if r.Code != 0 {
				return fmt.Errorf("API returned code %d", r.Code)
			}
		}

		return resourceDevice42DeviceRead(d, m)
	} else {
		return fmt.Errorf("incorrect response to query")
	}

}

func resourceDevice42SubnetRead(d *schema.ResourceData, m interface{}) error {

	log.Printf("[DEBUG] resourceDevice42SubnetRead - Starting reading using API for id %s", d.Id())
	resp, err := apiDevice42Get(m.(*resty.Client), fmt.Sprintf("/1.0/subnets/%s/", d.Id()), apiSubnetReadResponse{})

	if err != nil {
		log.Printf("[WARN] No subnet found: %s", d.Id())
		d.SetId("")
		return err
	}

	r := resp.Result().(*apiSubnetReadResponse)
	str := fmt.Sprintf("%v", r)
	log.Printf("[DEBUG] resourceDevice42SubnetRead - API data %s", str)
	d.Set("name", r.Name)
	d.Set("network", r.Network)
	d.Set("mask_bits", r.MaskBits)
	d.Set("vrf_group", r.VrfGroupName)
	d.Set("vrf_group_id", r.VrfGroupId)
	d.Set("parent_subnet_id", r.ParentSubnetId)
	d.Set("description", r.Description)
	d.Set("customer", r.Customer)
	fields := flattenCustomFields(r.CustomFields)
	d.Set("custom_fields", fields)
	return nil
}

func resourceDevice42SubnetUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	url := "/1.0/subnets/"
	device42SubnetUpdateFormData := map[string]string{}
	device42SubnetUpdateFormData["id"] = d.Id()
	if d.HasChange("name") {
		device42SubnetUpdateFormData["name"] = d.Get("name").(string)
	}
	if d.HasChange("mask_bits") {
		device42SubnetUpdateFormData["mask_bits"] = strconv.Itoa(d.Get("mask_bits").(int))
	}
	if d.HasChange("vrf_group") {
		device42SubnetUpdateFormData["vrf_group"] = d.Get("vrf_group").(string)
	}
	if d.HasChange("gateway") {
		device42SubnetUpdateFormData["gateway"] = d.Get("gateway").(string)
	}
	if d.HasChange("parent_subnet_id") {
		device42SubnetUpdateFormData["parent_subnet_id"] = strconv.Itoa(d.Get("parent_subnet_id").(int))
	}
	if d.HasChange("description") {
		device42SubnetUpdateFormData["description"] = d.Get("description").(string)
	}
	if d.HasChange("customer") {
		device42SubnetUpdateFormData["customer"] = d.Get("customer").(string)
	}
	log.Printf("[DEBUG] resourceDevice42SubnetUpdate - subnet : %s", d.Get("name").(string))

	resp, err := client.R().
		SetFormData(device42SubnetUpdateFormData).
		SetResult(apiResponse{}).
		Put(url)

	if d.HasChange("custom_fields") {
		updateList := setCustomFields(d)
		network := d.Get("network").(string)
		mask_bits := strconv.Itoa(d.Get("mask_bits").(int))
		vrf_group := d.Get("vrf_group").(string)
		for k, v := range updateList {
			formData := map[string]string{
				"network":   network,
				"mask_bits": mask_bits,
				"vrf_group": vrf_group,
				"key":       k,
				"value":     v.(string),
			}
			log.Printf("[DEBUG] resourceDevice42SubnetUpdate custom_fields  : %#v", formData)
			resp, err := client.R().
				SetFormData(formData).
				SetResult(apiResponse{}).
				Put("/1.0/custom_fields/subnet/")

			if err != nil {
				return err
			}

			r := resp.Result().(*apiResponse)
			log.Printf("[DEBUG] Result: %#v", r)
		}
	}

	if err != nil {
		return err
	}
	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] resourceDevice42SubnetUpdate - Result: %#v", r)

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
	log.Printf("[DEBUG] resourceDevice42SubnetDelete - Result: %#v", r)
	return nil
}
