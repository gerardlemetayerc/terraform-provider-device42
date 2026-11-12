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

// apiCustomerContactResponse is a single contact embedded in a parent customer's
// "Contacts" array, as returned by GET /1.0/customers/{id}/ and GET /1.0/customers/.
type apiCustomerContactResponse struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	ID      int64  `json:"id"`
}

// customerWithContacts is the parent customer shape including its "Contacts"
// array. apiCustomerReadResponse intentionally does not model the Contacts
// array, so this local wrapper carries them for contact lookup. The embedded
// apiCustomerReadResponse keeps the shared customer fields (Name, ID, ...)
// available as promoted fields.
type customerWithContacts struct {
	apiCustomerReadResponse
	Contacts []apiCustomerContactResponse `json:"Contacts"`
}

// apiCustomerContactsResponse is the wrapper returned by GET /1.0/customers/
// with the Contacts array modelled. Customer contacts have no list endpoint, so
// they are only readable inline inside a parent customer object.
type apiCustomerContactsResponse struct {
	Customers []customerWithContacts `json:"Customers"`
}

func resourceD42CustomerContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceDevice42CustomerContactCreate,
		Read:   resourceDevice42CustomerContactRead,
		Update: resourceDevice42CustomerContactUpdate,
		Delete: resourceDevice42CustomerContactDelete,

		Schema: map[string]*schema.Schema{
			"customer": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the parent customer the contact belongs to.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the customer contact.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the customer contact (e.g. Technical, Administrative).",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email address of the customer contact.",
			},
			"phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Phone number of the customer contact.",
			},
			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address of the customer contact.",
			},
			"contact_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the customer contact.",
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				client := m.(*resty.Client)
				id := d.Id()

				// Import id is encoded as "<customer_name>:<contact_name>".
				parts := strings.Split(id, ":")
				if len(parts) != 2 {
					return nil, fmt.Errorf("expected import id to be <customer_name>:<contact_name>, got %q", id)
				}

				customerName := parts[0]
				contactName := parts[1]

				contactId, err := getContactIdFromName(client, customerName, contactName)
				if err != nil {
					return nil, err
				}

				d.Set("customer", customerName)
				d.SetId(contactId)
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceDevice42CustomerContactCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	formData := map[string]string{}
	formData["customer"] = d.Get("customer").(string)
	formData["name"] = d.Get("name").(string)
	formData["type"] = d.Get("type").(string)

	if v := d.Get("email").(string); v != "" {
		formData["email"] = v
	}
	if v := d.Get("phone").(string); v != "" {
		formData["phone"] = v
	}
	if v := d.Get("address").(string); v != "" {
		formData["address"] = v
	}

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/customers/contacts/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42CustomerContactCreate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42CustomerContactCreate - Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42CustomerContactCreate - unable to extract customer contact id from response: %s", err)
	}

	d.SetId(strconv.Itoa(id))
	d.Set("contact_id", id)

	return resourceDevice42CustomerContactRead(d, m)
}

func resourceDevice42CustomerContactRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	customerName := d.Get("customer").(string)

	contactID, err := strconv.Atoi(d.Id())
	if err != nil {
		log.Printf("[WARN] Invalid customer contact id: %s", d.Id())
		d.SetId("")
		return nil
	}

	contact, err := getCustomerContact(client, customerName, contactID)
	if err != nil {
		log.Printf("[WARN] No customer contact found: %d", contactID)
		d.SetId("")
		return nil
	}

	d.Set("name", contact.Name)
	d.Set("type", contact.Type)
	d.Set("email", contact.Email)
	d.Set("phone", contact.Phone)
	d.Set("address", contact.Address)
	d.Set("contact_id", contact.ID)

	return nil
}

func resourceDevice42CustomerContactUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	formData := map[string]string{}
	formData["customer"] = d.Get("customer").(string)
	formData["name"] = d.Get("name").(string)
	formData["type"] = d.Get("type").(string)
	formData["contact_id"] = strconv.Itoa(d.Get("contact_id").(int))

	if v := d.Get("email").(string); v != "" {
		formData["email"] = v
	}
	if v := d.Get("phone").(string); v != "" {
		formData["phone"] = v
	}
	if v := d.Get("address").(string); v != "" {
		formData["address"] = v
	}

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/customers/contacts/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42CustomerContactUpdate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42CustomerContactUpdate - Result: %#v", r)

	return resourceDevice42CustomerContactRead(d, m)
}

func resourceDevice42CustomerContactDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting customer contact %s (ID: %s)", d.Get("name"), d.Id())

	url := fmt.Sprintf("/1.0/customers/contacts/%s/", d.Id())

	resp, err := client.R().
		SetResult(apiDeleteResponse{}).
		Delete(url)

	if err != nil {
		return err
	}

	r := resp.Result().(*apiDeleteResponse)
	log.Printf("[DEBUG] resourceDevice42CustomerContactDelete - Result: %#v", r)
	return nil
}

// getCustomerContact resolves a single contact by id under the customer whose
// name exactly matches customerName. Customer contacts have no list endpoint, so
// it fetches all customers and filters the parent client-side before scanning
// its Contacts array.
func getCustomerContact(client *resty.Client, customerName string, contactID int) (apiCustomerContactResponse, error) {
	var lookup apiCustomerContactsResponse
	url := fmt.Sprintf("/1.0/customers/")
	log.Printf("[DEBUG] getCustomerContact - Query customers on URL : %s", url)
	_, err := client.R().
		SetResult(&lookup).
		Get(url)

	if err != nil {
		return apiCustomerContactResponse{}, fmt.Errorf("error querying Device42 API: %s", err)
	}

	parent, err := filterCustomerByName(toPlainCustomers(lookup.Customers), customerName)
	if err != nil {
		return apiCustomerContactResponse{}, err
	}

	for _, c := range lookup.Customers {
		if c.ID == parent.ID {
			contact, ok := findContact(c, contactID)
			if ok {
				return contact, nil
			}
			return apiCustomerContactResponse{}, fmt.Errorf("no customer contact found with id: %d", contactID)
		}
	}

	return apiCustomerContactResponse{}, fmt.Errorf("no customer found with name: %s", customerName)
}

// findContact returns the contact whose id matches contactID within the parent
// customer's Contacts array, and whether such a contact exists.
func findContact(customer customerWithContacts, contactID int) (apiCustomerContactResponse, bool) {
	for _, contact := range customer.Contacts {
		if int(contact.ID) == contactID {
			return contact, true
		}
	}
	return apiCustomerContactResponse{}, false
}

// findContactByName returns all contacts whose name matches contactName within
// the parent customer's Contacts array. It errors when no contact matches, so the
// caller can detect an ambiguous (duplicate) name by checking the slice length.
func findContactByName(customer customerWithContacts, contactName string) ([]apiCustomerContactResponse, error) {
	var matches []apiCustomerContactResponse
	for _, contact := range customer.Contacts {
		if contact.Name == contactName {
			matches = append(matches, contact)
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no customer contact found with name: %s", contactName)
	}

	return matches, nil
}

// getContactIdFromName resolves a contact id by fetching all customers, filtering
// the parent customer by name client-side, then finding the contact by name in
// its Contacts array. It errors on ambiguity (duplicate contact names).
func getContactIdFromName(client *resty.Client, customerName, contactName string) (string, error) {
	var lookup apiCustomerContactsResponse
	url := fmt.Sprintf("/1.0/customers/")
	log.Printf("[DEBUG] getContactIdFromName - Query customers on URL : %s", url)
	_, err := client.R().
		SetResult(&lookup).
		Get(url)

	log.Printf("[DEBUG] Data from customer lookup : %#v", lookup)
	if err != nil {
		return "", fmt.Errorf("error querying Device42 API: %s", err)
	}

	parent, err := filterCustomerByName(toPlainCustomers(lookup.Customers), customerName)
	if err != nil {
		return "", err
	}

	for _, c := range lookup.Customers {
		if c.ID == parent.ID {
			matches, err := findContactByName(c, contactName)
			if err != nil {
				return "", err
			}

			if len(matches) > 1 {
				return "", fmt.Errorf("more than one customer contact found with name: %s", contactName)
			}

			return strconv.FormatInt(matches[0].ID, 10), nil
		}
	}

	return "", fmt.Errorf("no customer found with name: %s", customerName)
}

// toPlainCustomers extracts the shared apiCustomerReadResponse fields from the
// richer customerWithContacts shape so the shared filterCustomerByName helper can
// be reused for client-side by-name filtering.
func toPlainCustomers(customers []customerWithContacts) []apiCustomerReadResponse {
	plain := make([]apiCustomerReadResponse, len(customers))
	for i, c := range customers {
		plain[i] = c.apiCustomerReadResponse
	}
	return plain
}
