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
	HostEnv string = "D42_HOST"
	// Environment variable to configure the device42 api username attribute
	UsernameEnv string = "D42_USER"
	// Environment variable to configure the device42 api password attribute
	PasswordEnv string = "D42_PASS"
)

// Provider -- main device42 provider structure
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			// -- API Interaction Definitions --
			"D42_HOST": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc(
					HostEnv,
					"",
				),
				Description: "The device42 server to interact with.",
			},
			"D42_PASSWORD": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc(
					PasswordEnv,
					"",
				),
				Description: "The password to authenticate with Device42.",
			},
			"D42_USERNAME": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc(
					UsernameEnv,
					"",
				),
				Description: "The username to authenticate with Device42.",
			},
			"D42_TLS_INSECURE": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "Whether to perform TLS cert verification on the server's certificate. " +
					"Defaults to `false`.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"device42_device": resourceD42Device(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	host := d.Get("D42_HOST").(string)
	username := d.Get("D42_USERNAME").(string)
	password := d.Get("D42_PASSWORD").(string)
	tlsInsecure := d.Get("D42_TLS_INSECURE").(bool)

	if host == "" {
		return nil, fmt.Errorf("no Device42 host was provided")
	}

	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: tlsInsecure})
	client.SetBaseURL(fmt.Sprintf("https://%s/api", host))
	client.SetBasicAuth(username, password)

	return client, nil
}
