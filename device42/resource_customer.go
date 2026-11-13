package device42

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// apiCustomerReadResponse is the shape returned by GET /1.0/customers/{id}/ which
// wraps the customer under a "Customers" key. Only the fields this resource
// models are present; the "Contacts" array is intentionally not modelled.
type apiCustomerReadResponse struct {
	Name         string        `json:"name"`
	Manager      string        `json:"manager"`
	Notes        string        `json:"notes"`
	ContactInfo  string        `json:"contact_info"`
	Tags         []interface{} `json:"tags"`
	CustomFields []customField `json:"custom_fields"`
	ID           int64         `json:"id"`
}

// apiCustomersResponse is the wrapper returned by GET /1.0/customers/{id}/ and
// GET /1.0/customers/. When the customer does not exist the API returns HTTP 200
// with an empty array under "Customers" instead of an error.
type apiCustomersResponse struct {
	Customers []apiCustomerReadResponse `json:"Customers"`
}

func resourceD42Customer() *schema.Resource {
	return &schema.Resource{
		Create: resourceDevice42CustomerCreate,
		Read:   resourceDevice42CustomerRead,
		Update: resourceDevice42CustomerUpdate,
		Delete: resourceDevice42CustomerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the customer.",
			},
			"customer_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the customer.",
			},
			"manager": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Manager of the customer.",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Any additional notes.",
			},
			"contact_info": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Contact information for the customer.",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Add or update tags to a customer (comma separated).",
			},
			"custom_fields": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				Description:      "Any custom fields that will be used in device42.",
				DiffSuppressFunc: suppressCustomFieldsDiffs,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				client := m.(*resty.Client)
				name := d.Id()

				customerId, err := getCustomerIdFromName(client, name)
				if err != nil {
					return nil, err
				}

				d.SetId(customerId)
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceDevice42CustomerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	name := d.Get("name").(string)
	formData := map[string]string{}
	formData["name"] = name

	if v := d.Get("manager").(string); v != "" {
		formData["manager"] = v
	}
	if v := d.Get("notes").(string); v != "" {
		formData["notes"] = v
	}
	if v := d.Get("contact_info").(string); v != "" {
		formData["contact_info"] = v
	}
	if v := d.Get("tags").(string); v != "" {
		formData["tags"] = v
	}

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/customers/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42CustomerCreate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42CustomerCreate - Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42CustomerCreate - unable to extract customer id from response: %s", err)
	}

	d.SetId(strconv.Itoa(id))
	d.Set("customer_id", id)

	if d.Get("custom_fields") != nil {
		fields := d.Get("custom_fields").(map[string]interface{})
		bulkFields := []string{}

		for k, v := range fields {
			bulkFields = append(bulkFields, fmt.Sprintf("%v:%v", k, v))
		}

		log.Printf("[DEBUG] resourceDevice42CustomerCreate - Pushing custom_fields %s", bulkFields)
		resp, err := client.R().
			SetFormData(map[string]string{
				"name":        name,
				"bulk_fields": strings.Join(bulkFields, ","),
			}).
			SetResult(apiResponse{}).
			Put("/1.0/custom_fields/customer/")

		if err != nil {
			return err
		}

		r := resp.Result().(*apiResponse)

		if r.Code != 0 {
			return fmt.Errorf("resourceDevice42CustomerCreate - api returned code %d", r.Code)
		}
	}

	return resourceDevice42CustomerRead(d, m)
}

func resourceDevice42CustomerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	d42Url := fmt.Sprintf("/1.0/customers/%s/", d.Id())
	log.Printf("[DEBUG] resourceDevice42CustomerRead - reading customer on target URL: %s", d42Url)

	var resp apiCustomersResponse
	_, err := client.R().
		SetResult(&resp).
		Get(d42Url)

	if err != nil {
		log.Printf("[WARN] No customer found: %s", d.Id())
		d.SetId("")
		return err
	}

	if len(resp.Customers) == 0 {
		log.Printf("[WARN] No customer found: %s", d.Id())
		d.SetId("")
		return nil
	}

	c := resp.Customers[0]
	d.Set("name", c.Name)
	d.Set("customer_id", c.ID)
	d.Set("manager", c.Manager)
	d.Set("notes", c.Notes)
	d.Set("contact_info", c.ContactInfo)
	d.Set("tags", joinBuildingList(c.Tags))
	d.Set("custom_fields", flattenCustomFields(c.CustomFields))

	return nil
}

func resourceDevice42CustomerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	name := d.Get("name").(string)
	formData := map[string]string{}
	formData["name"] = name

	if d.HasChange("name") {
		// name is ForceNew so this should not occur, but guard anyway.
		formData["new_name"] = name
	}
	if d.HasChange("manager") {
		formData["manager"] = d.Get("manager").(string)
	}
	if d.HasChange("notes") {
		formData["notes"] = d.Get("notes").(string)
	}
	if d.HasChange("contact_info") {
		formData["contact_info"] = d.Get("contact_info").(string)
	}
	if d.HasChange("tags") {
		formData["tags"] = d.Get("tags").(string)
	}
	log.Printf("[DEBUG] resourceDevice42CustomerUpdate - Pushing new informations : %#v", formData)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/customers/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42CustomerUpdate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42CustomerUpdate - Result: %#v", r)

	if d.HasChange("custom_fields") {
		updateList := setCustomFields(d)
		for k, v := range updateList {
			log.Printf("[DEBUG] resourceDevice42CustomerUpdate - Pushing custom field %s : %v", k, v)
			resp, err := client.R().
				SetFormData(map[string]string{
					"name":  name,
					"key":   k,
					"value": v.(string),
				}).
				SetResult(apiResponse{}).
				Put("/1.0/custom_fields/customer/")

			if err != nil {
				return err
			}

			r := resp.Result().(*apiResponse)
			log.Printf("[DEBUG] resourceDevice42CustomerUpdate - Result: %#v", r)
		}
	}

	return resourceDevice42CustomerRead(d, m)
}

func resourceDevice42CustomerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting customer %s (ID: %s)", d.Get("name"), d.Id())

	url := fmt.Sprintf("/1.0/customers/%s/", d.Id())

	resp, err := client.R().
		SetResult(apiDeleteResponse{}).
		Delete(url)

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42CustomerDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiDeleteResponse)
	log.Printf("[DEBUG] resourceDevice42CustomerDelete - Result: %#v", r)
	return nil
}

// filterCustomerByName returns the single customer whose name exactly matches
// the given name. Because GET /1.0/customers/ ignores the ?name= query param and
// returns ALL customers, the by-name lookup must be done client-side, erroring
// when no match or more than one match is found.
func filterCustomerByName(customers []apiCustomerReadResponse, name string) (apiCustomerReadResponse, error) {
	var matches []apiCustomerReadResponse
	for _, c := range customers {
		if c.Name == name {
			matches = append(matches, c)
		}
	}

	if len(matches) == 0 {
		return apiCustomerReadResponse{}, fmt.Errorf("no customer found with name: %s", name)
	}

	if len(matches) > 1 {
		return apiCustomerReadResponse{}, fmt.Errorf("more than one customer found with name: %s", name)
	}

	return matches[0], nil
}

func getCustomerIdFromName(client *resty.Client, name string) (string, error) {
	var lookup apiCustomersResponse
	url := fmt.Sprintf("/1.0/customers/")
	log.Printf("[DEBUG] Query for customer name on URL : %s", url)
	_, err := client.R().
		SetResult(&lookup).
		Get(url)

	log.Printf("[DEBUG] Data from customer lookup : %#v", lookup)
	if err != nil {
		return "", fmt.Errorf("error querying Device42 API: %s", err)
	}

	c, err := filterCustomerByName(lookup.Customers, name)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(c.ID, 10), nil
}
