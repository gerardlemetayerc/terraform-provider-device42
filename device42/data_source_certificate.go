package device42

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceD42Certificate() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42CertificateRead,
		Description: "Read a Device42 certificate.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Certificate common name / subject CN, matches issued_to.",
			},
			"certificate_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the certificate.",
			},
			"issued_to": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"serial_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"san": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_from": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_to": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signature_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signature_hash": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extended_key_usage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issued_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_cert": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"days_to_expiry": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"key_agreement_usage": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"encipher_only_usage": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"key_encipherment_usage": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"data_encipherment_usage": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"crl_sign_usage": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"decipher_only_usage": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"content_commitment_usage": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"digital_signature_usage": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"key_cert_sign_usage": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"custom_fields": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func datasourceD42CertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)

	name := d.Get("name").(string)
	certificateID := d.Get("certificate_id").(int)

	if name != "" {
		// GET /1.0/certificates/ returns ALL certificates, so the by-name lookup
		// must be filtered client-side.
		url := fmt.Sprintf("/1.0/certificates/")
		log.Printf("[DEBUG] datasourceD42CertificateRead - Query certificates on URL : %s", url)

		var lookup apiCertificateListResponse
		_, err := client.R().
			SetResult(&lookup).
			Get(url)

		if err != nil {
			log.Printf("[WARN] No certificate found: %s", name)
			d.SetId("")
			return nil
		}

		c, err := filterCertificateByName(lookup.CertificateDetails, name)
		if err != nil {
			log.Printf("[ERROR] %s", err)
			d.SetId("")
			return nil
		}

		d.SetId(strconv.FormatInt(c.ID, 10))
		d.Set("certificate_id", c.ID)
		setCertificateFromDetails(d, c)
		return nil
	}

	if certificateID > 0 {
		url := fmt.Sprintf("/1.0/certificates/%d/", certificateID)
		log.Printf("[DEBUG] datasourceD42CertificateRead - Query certificate on URL : %s", url)

		var result apiCertificateResponse
		apiResp, err := client.R().
			SetResult(&result).
			Get(url)

		if err != nil {
			log.Printf("[WARN] No certificate found: %d", certificateID)
			d.SetId("")
			return nil
		}

		// resty does not surface a 404 as an error, so a deleted/unknown
		// certificate is detected via the HTTP status code.
		if apiResp.StatusCode() == http.StatusNotFound || result.CertificateDetails.ID == 0 {
			log.Printf("[WARN] No certificate found: %d", certificateID)
			d.SetId("")
			return nil
		}

		d.SetId(strconv.Itoa(certificateID))
		d.Set("certificate_id", result.CertificateDetails.ID)
		setCertificateFromDetails(d, result.CertificateDetails)
		return nil
	}

	return fmt.Errorf("either 'name' or 'certificate_id' must be set to look up a Device42 certificate")
}
