package device42

import (
	"crypto/tls"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Environment variables the provider recognizes for configuration
const (
	// Environment variable to configure the device42 api host
	HostEnv string = "d42_host"
	// Environment variable to configure the device42 api username attribute
	UsernameEnv string = "d42_username"
	// Environment variable to configure the device42 api password attribute
	PasswordEnv string = "d42_password"
)

// Provider -- main device42 provider structure
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			// -- API Interaction Definitions --
			"d42_host": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc(
					HostEnv,
					"",
				),
				Description: "The device42 server to interact with.",
			},
			"d42_password": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc(
					PasswordEnv,
					"",
				),
				Description: "The password to authenticate with Device42.",
			},
			"d42_username": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc(
					UsernameEnv,
					"",
				),
				Description: "The username to authenticate with Device42.",
			},
			"d42_tls_unsecure": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "Whether to perform TLS cert verification on the server's certificate. " +
					"Defaults to `false`.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"device42_device": resourceD42Device(),
			"device42_subnet": resourceD42Subnet(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	host := d.Get("d42_host").(string)
	username := d.Get("d42_username").(string)
	password := d.Get("d42_password").(string)
	tlsInsecure := d.Get("d42_tls_unsecure").(bool)

	if host == "" {
		return nil, fmt.Errorf("no Device42 host was provided")
	}

	if username == "" {
		return nil, fmt.Errorf("no username was provided")
	}

	if password == "" {
		return nil, fmt.Errorf("no password was provided")
	}

	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: tlsInsecure})
	client.SetBaseURL(fmt.Sprintf("https://%s/api", host))
	client.SetBasicAuth(username, password)

	return client, nil
}
