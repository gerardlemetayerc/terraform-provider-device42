package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiBusinessAppsElement struct {
	BusinessAppId   int    `json:"business_app_id"`
	BusinessAppName int    `json:"business_app_name"`
	Element         string `json:"element"`
}

type apiBusinessAppsElementApiResponse struct {
	BusinessappElements []apiBusinessAppsElement `json:"businessapp_elements"`
}

func resourceD42BusinessAppsElement() *schema.Resource {
	return &schema.Resource{
		Description: "device42_businessappelement can be use to manage Business Applications element",
		Create:      resourceDevice42BusinessAppsElementCreate,
		Update:      resourceDevice42BusinessAppsElementCreate,
		Read:        resourceDevice42BusinessAppsElementRead,
		Delete:      resourceDevice42BusinessAppsElementDelete,
		Schema: map[string]*schema.Schema{
			"businessapp_id": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: "The ID of an existing Business Application to add elements (devices) to.				.",
			},
			"device_id": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: "ID of an element (device) to add to the business app.",
			},
		},
	}
}

func resourceDevice42BusinessAppsElementRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] resourceDevice42BusinessAppsElementRead - Starting reading using API for id %s", d.Id())
	resp, err := client.R().
		SetResult(apiBusinessAppsElementApiResponse{}).
		Get(fmt.Sprintf("/1.0/businessapps/elements/%s/", d.Id()))

	if err != nil {
		log.Printf("[WARN] No data found for id %s", d.Id())
		return err
	}

	r := resp.Result().(*apiBusinessAppsElementApiResponse)
	str := fmt.Sprintf("%v", r)
	log.Printf("[DEBUG] resourceDevice42BusinessAppsElementRead - API data %s", str)
	d.Set("businessapp_id", r.BusinessappElements[0].BusinessAppId)
	d.Set("device_id", r.BusinessappElements[0].Element)
	return nil
}

func resourceDevice42BusinessAppsElementCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	businessapp_id := d.Get("businessapp_id").(int)
	device_id := d.Get("device_id").(int)

	formData := map[string]string{
		"businessapp_id": strconv.Itoa(businessapp_id),
		"device_id":      strconv.Itoa(device_id),
	}

	log.Printf("[DEBUG] resourceDevice42BusinessAppsElementCreate - Starting reading using API for id %s", d.Id())
	resp, err := client.R().
		SetResult(apiReadData{}).
		SetFormData(formData).
		Post("/1.0/businessapps/elements/")

	if err != nil {
		log.Printf("issue during creation. Error code : %s", d.Id())
		return err
	}

	r := resp.Result().(*apiResponse)

	if len(r.Msg) < 1 {
		str := fmt.Sprintf("%v", r)
		return fmt.Errorf("please check account permission or credentials - api returned :  %s", str)
	}

	log.Printf("[DEBUG] Result: %#v", r)
	id := int(r.Msg[1].(float64))

	// Set ID after Business App creation
	d.SetId(strconv.Itoa(id))
	return resourceDevice42BusinessAppsElementRead(d, m)
}

func resourceDevice42BusinessAppsElementDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	businessapp_id := d.Get("businessapp_id").(int)
	device_id := d.Get("device_id").(int)

	formData := map[string]string{
		"businessapp_id": strconv.Itoa(businessapp_id),
		"device_id":      strconv.Itoa(device_id),
	}

	log.Printf("Deleting Business Apps Element")
	resp, err := client.R().
		SetResult(apiResponse{}).
		SetFormData(formData).
		Delete("/1.0/businessapps/elements/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)
	log.Printf("[DEBUG] Result: %#v", r)
	return nil
}
