package device42

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42Customer() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42CustomerRead,
		Description: "Read a Device42 customer.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the customer.",
			},
			"customer_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the customer.",
			},
			"manager": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"contact_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
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

func datasourceD42CustomerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	customerID := d.Get("customer_id").(int)

	if name != "" {
		// GET /1.0/customers/ ignores the ?name= query param and returns ALL
		// customers, so the by-name lookup must be filtered client-side.
		url := fmt.Sprintf("/1.0/customers/")
		log.Printf("[DEBUG] datasourceD42CustomerRead - Query customers on URL : %s", url)

		var lookup apiCustomersResponse
		apiResp, err := client.R().
			SetResult(&lookup).
			Get(url)

		if err != nil {
			log.Printf("[WARN] No customer found: %s", name)
			d.SetId("")
			return err
		}
		if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != http.StatusNotFound {
			return fmt.Errorf("datasourceD42CustomerRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
		}

		c, err := filterCustomerByName(lookup.Customers, name)
		if err != nil {
			log.Printf("[ERROR] %s", err)
			d.SetId("")
			return nil
		}

		d.SetId(strconv.FormatInt(c.ID, 10))
		d.Set("customer_id", c.ID)
		d.Set("name", c.Name)
		d.Set("manager", c.Manager)
		d.Set("notes", c.Notes)
		d.Set("contact_info", c.ContactInfo)
		d.Set("tags", joinBuildingList(c.Tags))
		d.Set("custom_fields", flattenCustomFields(c.CustomFields))
		return nil
	}

	if customerID > 0 {
		url := fmt.Sprintf("/1.0/customers/%d/", customerID)
		log.Printf("[DEBUG] datasourceD42CustomerRead - Query customer on URL : %s", url)

		var resp apiCustomersResponse
		apiResp, err := client.R().
			SetResult(&resp).
			Get(url)

		if err != nil {
			log.Printf("[WARN] No customer found: %d", customerID)
			d.SetId("")
			return err
		}
		if apiResp.StatusCode() == http.StatusNotFound {
			log.Printf("[WARN] No customer found: %d", customerID)
			d.SetId("")
			return nil
		}
		if apiResp.StatusCode() >= 400 {
			return fmt.Errorf("datasourceD42CustomerRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
		}

		if len(resp.Customers) == 0 {
			log.Printf("[ERROR] No customer found with id: %d", customerID)
			d.SetId("")
			return nil
		}

		c := resp.Customers[0]
		d.SetId(strconv.Itoa(customerID))
		d.Set("customer_id", c.ID)
		d.Set("name", c.Name)
		d.Set("manager", c.Manager)
		d.Set("notes", c.Notes)
		d.Set("contact_info", c.ContactInfo)
		d.Set("tags", joinBuildingList(c.Tags))
		d.Set("custom_fields", flattenCustomFields(c.CustomFields))
		return nil
	}

	return fmt.Errorf("either 'name' or 'customer_id' must be set to look up a Device42 customer")
}
