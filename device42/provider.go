package device42

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Environment variables the provider recognizes for configuration
const (
	// Environment variable to configure the device42 api host
	HostEnv string = "D42_HOST"
	// Environment variable to configure the device42 api username attribute
	UsernameEnv string = "D42_USERNAME"
	// Environment variable to configure the device42 api password attribute
	PasswordEnv string = "D42_PASSWORD"
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
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
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
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Timeout in seconds for API requests.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"device42_device":               resourceD42Device(),
			"device42_subnet":               resourceD42Subnet(),
			"device42_subnet_suggested":     resourceD42SubnetSuggested(),
			"device42_businessapp":          resourceD42BusinessApps(),
			"device42_businessappelement":   resourceD42BusinessAppsElement(),
			"device42_vlan":                 resourceD42Vlans(),
			"device42_ip":                   resourceD42Ip(),
			"device42_building":             resourceD42Building(),
			"device42_room":                 resourceD42Room(),
			"device42_rack":                 resourceD42Rack(),
			"device42_customer":             resourceD42Customer(),
			"device42_customer_contact":     resourceD42CustomerContact(),
			"device42_certificate":          resourceD42Certificate(),
			"device42_vendor":               resourceD42Vendor(),
			"device42_software":             resourceD42Software(),
			"device42_software_license_key": resourceD42SoftwareLicenseKey(),
			"device42_cable":                resourceD42Cable(),
			"device42_power_circuit":        resourceD42PowerCircuit(),
			"device42_pdu_model":            resourceD42PduModel(),
			"device42_pdu":                  resourceD42Pdu(),
			"device42_vrf_group":            resourceD42VrfGroup(),
			"device42_subnet_category":      resourceD42SubnetCategory(),
			"device42_mac":                  resourceD42Mac(),
			"device42_ip_nat":               resourceD42IpNat(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"device42_device":           datasourceD42Device(),
			"device42_subnet":           datasourceD42Subnet(),
			"device42_subnets":          datasourceD42Subnets(),
			"device42_ip":               datasourceD42Ip(),
			"device42_ips":              datasourceD42Ips(),
			"device42_vlan":             datasourceD42Vlan(),
			"device42_vlans":            datasourceD42Vlans(),
			"device42_suggestedIp":      datasourceD42SuggestedIp(),
			"device42_businessapp":      datasourceD42BusinessApp(),
			"device42_password":         datasourceD42Password(),
			"device42_building":         datasourceD42Building(),
			"device42_room":             datasourceD42Room(),
			"device42_rack":             datasourceD42Rack(),
			"device42_customer":         datasourceD42Customer(),
			"device42_customer_contact": datasourceD42CustomerContact(),
			"device42_certificate":      datasourceD42Certificate(),
			"device42_service_level":    datasourceD42ServiceLevel(),
			"device42_vendor":           datasourceD42Vendor(),
			"device42_software":         datasourceD42Software(),
			"device42_software_details": datasourceD42SoftwareDetails(),
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

	timeout := d.Get("timeout").(int)

	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: tlsInsecure})
	client.SetBaseURL(fmt.Sprintf("https://%s/api", host))
	client.SetBasicAuth(username, password)
	client.SetTimeout(time.Duration(timeout) * time.Second)

	return client, nil
}
