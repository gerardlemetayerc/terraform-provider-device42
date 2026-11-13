package device42

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42CustomerContact() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42CustomerContactRead,
		Description: "Read a Device42 customer contact.",
		Schema: map[string]*schema.Schema{
			"customer": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the parent customer the contact belongs to.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the customer contact.",
			},
			"contact_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the customer contact.",
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"phone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceD42CustomerContactRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	customerName := d.Get("customer").(string)
	contactName := d.Get("name").(string)
	contactID := d.Get("contact_id").(int)

	if customerName == "" {
		return fmt.Errorf("'customer' must be set to look up a Device42 customer contact")
	}

	if contactName == "" && contactID == 0 {
		return fmt.Errorf("either 'name' or 'contact_id' must be set to look up a Device42 customer contact")
	}

	// GET /1.0/customers/ ignores the ?name= query param and returns ALL
	// customers, so the by-name lookup must be filtered client-side.
	url := fmt.Sprintf("/1.0/customers/")
	log.Printf("[DEBUG] datasourceD42CustomerContactRead - Query customers on URL : %s", url)

	var lookup apiCustomerContactsResponse
	apiResp, err := client.R().
		SetResult(&lookup).
		Get(url)

	if err != nil {
		log.Printf("[WARN] No customer found: %s", customerName)
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() >= 400 && apiResp.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("datasourceD42CustomerContactRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	parent, err := filterCustomerByName(toPlainCustomers(lookup.Customers), customerName)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		d.SetId("")
		return nil
	}

	for _, c := range lookup.Customers {
		if c.ID == parent.ID {
			var contact apiCustomerContactResponse
			if contactID > 0 {
				var ok bool
				contact, ok = findContact(c, contactID)
				if !ok {
					log.Printf("[WARN] No customer contact found: %d", contactID)
					d.SetId("")
					return nil
				}
			} else {
				matches, err := findContactByName(c, contactName)
				if err != nil {
					log.Printf("[ERROR] %s", err)
					d.SetId("")
					return nil
				}

				if len(matches) > 1 {
					log.Printf("[ERROR] more than one customer contact found with name: %s", contactName)
					d.SetId("")
					return nil
				}

				contact = matches[0]
			}

			d.SetId(strconv.FormatInt(contact.ID, 10))
			d.Set("contact_id", contact.ID)
			d.Set("name", contact.Name)
			d.Set("type", contact.Type)
			d.Set("email", contact.Email)
			d.Set("phone", contact.Phone)
			d.Set("address", contact.Address)
			return nil
		}
	}

	log.Printf("[WARN] No customer found: %s", customerName)
	d.SetId("")
	return nil
}
