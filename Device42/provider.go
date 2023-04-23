package device42

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
