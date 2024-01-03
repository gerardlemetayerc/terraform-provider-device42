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
	
	if v, ok := d.GetOk("label"); ok {
		queryParams.Set("label", v.(string))
	}
	if v, ok := d.GetOk("username"); ok {
		queryParams.Set("username", v.(string))
	}
	if v, ok := d.GetOk("appcomp"); ok {
		queryParams.Set("appcomp", v.(string))
	}
	if v, ok := d.GetOk("id"); ok {
		queryParams.Set("id", strconv.Itoa(v.(int)))
	}

	var resp datasourceD42PasswordResponse
	_, err := client.R().
		SetResult(resp{}).
		SetQueryParamsFromValues(queryParams).
		Get("/api/1.0/passwords/")

	if err != nil {
		log.Printf("[WARN] No password found: %s", err)
		return err
	}

	if len(resp.Passwords) > 0 {
		d.SetId(strconv.Itoa(int((resp.Passwords[0]).ID)))
		d.Set("username", resp.Passwords[0].Username)
		d.Set("device", resp.Passwords[0].Device)
		d.Set("label", resp.Passwords[0].Label)
		d.Set("password", resp.Passwords[0].Password)
	}

	return nil
}
