package device42

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

const defaultBaseURL = "https://device42.example.com"

// Provider retourne une instance du provider Terraform pour Device42.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     defaultBaseURL,
				Description: "URL de base de l'API Device42.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Nom d'utilisateur pour l'authentification.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Mot de passe pour l'authentification.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"device42_virtual_machine": resourceVirtualMachine(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure configure le client Device42 à partir des paramètres du provider.
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	baseURL := d.Get("base_url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	client := NewClient(baseURL, httpClient)

	// Authentification auprès de Device42.
	authErr := client.Authenticate(username, password)
	if authErr != nil {
		return nil, diag.FromErr(authErr)
	}

	return client, nil
}

// providerError renvoie une erreur formatée pour les erreurs spécifiques au provider.
func providerError(msg string, args ...interface{}) error {
	return fmt.Errorf("device42: "+msg, args...)
}

func resourceDevice42VirtualMachine() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceDevice42VirtualMachineCreate,
        ReadContext:   resourceDevice42VirtualMachineRead,
        UpdateContext: resourceDevice42VirtualMachineUpdate,
        DeleteContext: resourceDevice42VirtualMachineDelete,

        Schema: map[string]*schema.Schema{
            "name": {
                Type:     schema.TypeString,
                Required: true,
            },
            "ip_address": {
                Type:     schema.TypeString,
                Required: true,
            },
            "subnet": {
                Type:     schema.TypeString,
                Required: true,
            },
            "application": {
                Type:     schema.TypeString,
                Required: true,
            },
        },
    }
}

func resourceDevice42VirtualMachineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*device42.Client)
    name := d.Get("name").(string)
    ipAddress := d.Get("ip_address").(string)
    subnet := d.Get("subnet").(string)
    application := d.Get("application").(string)

    vm, err := client.CreateVirtualMachine(name, ipAddress, subnet, application)
    if err != nil {
        return diag.FromErr(err)
    }

    d.SetId(vm.ID)
    return nil
}

func resourceDevice42VirtualMachineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*device42.Client)
    vm, err := client.GetVirtualMachine(d.Id())
    if err != nil {
        return diag.FromErr(err)
    }

    d.Set("name", vm.Name)
    d.Set("ip_address", vm.IPAddress)
    d.Set("subnet", vm.Subnet)
    d.Set("application", vm.Application)

    return nil
}

func resourceDevice42VirtualMachineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*device42.Client)
    name := d.Get("name").(string)
    ipAddress := d.Get("ip_address").(string)
    subnet := d.Get("subnet").(string)
    application := d.Get("application").(string)

    _, err := client.UpdateVirtualMachine(d.Id(), name, ipAddress, subnet, application)
    if err != nil {
        return diag.FromErr(err)
    }

    return nil
}

func resourceDevice42VirtualMachineDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*device42.Client)

    err := client.DeleteVirtualMachine(d.Id())
    if err != nil {
        return diag.FromErr(err)
    }

    d.SetId("")
    return nil
}
