package device42

import (
	"log"
	"net/url"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type datasourceD42PasswordResponse struct {
	Passwords []struct {
		Username     string `json:"username"`
		LastPWChange string `json:"last_pw_change"`
		Notes        string `json:"notes"`
		Label        string `json:"label"`
		FirstAdded   string `json:"first_added"`
		Password     string `json:"password"`
		ID           int    `json:"id"`
	} `json:"Passwords"`
}

func datasourceD42Password() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42PasswordRead,
		Description: "Retrieve passwords information.",
		Schema: map[string]*schema.Schema{
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Category of the password.",
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Label of the password.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Username associated with the password.",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Password of the account.",
			},
			"device": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Device associated with the password.",
			},
			"appcomp": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Application component associated with the password.",
			},
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the password.",
			},
		},
	}
}

func datasourceD42PasswordRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	queryParams := make(url.Values)

	queryParams.Set("plain_text", "yes")

	if v, ok := d.GetOk("category"); ok {
		queryParams.Set("category", v.(string))
	}
	if v, ok := d.GetOk("label"); ok {
		queryParams.Set("label", v.(string))
	}
	if v, ok := d.GetOk("username"); ok {
		queryParams.Set("username", v.(string))
	}
	if v, ok := d.GetOk("device"); ok {
		queryParams.Set("device", v.(string))
	}
	if v, ok := d.GetOk("appcomp"); ok {
		queryParams.Set("appcomp", v.(string))
	}
	if v, ok := d.GetOk("id"); ok {
		queryParams.Set("id", strconv.Itoa(v.(int)))
	}

	resp, err := client.R().
		SetResult(datasourceD42PasswordResponse{}).
		SetQueryParamsFromValues(queryParams).
		Get("/api/1.0/passwords/")

	if err != nil {
		log.Printf("[WARN] No password found: %s", err)
		return err
	}

	r := resp.Result().(*datasourceD42PasswordResponse)
	if len(r.Passwords) > 0 {
		d.SetId(strconv.Itoa(int((r.Passwords[0]).id)))
		d.Set("device_id", (r.Passwords[0]).DeviceID)
		d.Set("category", r.Passwords[0].Category)
		d.Set("username", r.Passwords[0].Username)
		d.Set("device", r.Passwords[0].Device)
		d.Set("label", r.Passwords[0].Label)
		d.Set("password", r.Passwords[0].Password)
	}

	return nil
}
