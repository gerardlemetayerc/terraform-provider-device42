package device42

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// apiCertificateDetails is the shape of a single certificate object returned by
// /1.0/certificates/. The nullable fields (issued_to, issued_by, parent_cert,
// subject, san, vendor) are typed as string so JSON null decodes to "". The
// version field is tolerant of null via apiStringOrNumber; the bool usage flags
// are returned by the API and are modelled as read-only.
type apiCertificateDetails struct {
	ID                     int64             `json:"id"`
	IssuedTo               string            `json:"issued_to"`
	IssuedBy               string            `json:"issued_by"`
	ParentCert             string            `json:"parent_cert"`
	SerialNumber           string            `json:"serial_number"`
	Subject                string            `json:"subject"`
	SAN                    string            `json:"san"`
	ValidFrom              string            `json:"valid_from"`
	ValidTo                string            `json:"valid_to"`
	Vendor                 string            `json:"vendor"`
	SignatureAlgorithm     string            `json:"signature_algorithm"`
	SignatureHash          string            `json:"signature_hash"`
	ExtendedKeyUsage       string            `json:"extended_key_usage"`
	Version                apiStringOrNumber `json:"version"`
	DaysToExpiry           int               `json:"days_to_expiry"`
	CustomFields           []customField     `json:"custom_fields"`
	KeyAgreementUsage      bool              `json:"key_agreement_usage"`
	EncipherOnlyUsage      bool              `json:"encipher_only_usage"`
	KeyEnciphermentUsage   bool              `json:"key_encipherment_usage"`
	DataEnciphermentUsage  bool              `json:"data_encipherment_usage"`
	CRLSignUsage           bool              `json:"crl_sign_usage"`
	DecipherOnlyUsage      bool              `json:"decipher_only_usage"`
	ContentCommitmentUsage bool              `json:"content_commitment_usage"`
	DigitalSignatureUsage  bool              `json:"digital_signature_usage"`
	KeyCertSignUsage       bool              `json:"key_cert_sign_usage"`
}

// apiCertificateResponse is the shape returned by GET /1.0/certificates/{id}/
// which wraps the single certificate under "certificate_details" as a single
// OBJECT (not an array).
type apiCertificateResponse struct {
	CertificateDetails apiCertificateDetails `json:"certificate_details"`
}

// apiCertificateListResponse is the shape returned by GET /1.0/certificates/
// which lists certificates under "certificate_details" as an ARRAY.
type apiCertificateListResponse struct {
	TotalCount         int                     `json:"total_count"`
	CertificateDetails []apiCertificateDetails `json:"certificate_details"`
}

func resourceD42Certificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceDevice42CertificateCreate,
		Read:   resourceDevice42CertificateRead,
		Update: resourceDevice42CertificateUpdate,
		Delete: resourceDevice42CertificateDelete,

		Schema: map[string]*schema.Schema{
			"issued_to": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Certificate common name / subject CN (e.g. \"*.example.com\").",
			},
			"serial_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Serial number of the certificate.",
			},
			"subject": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subject of the certificate.",
			},
			"san": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subject alternative names of the certificate.",
			},
			"valid_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Certificate validity start date (YYYY-MM-DD).",
			},
			"valid_to": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Certificate validity end date (YYYY-MM-DD).",
			},
			"vendor": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Vendor of the certificate.",
			},
			"signature_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Signature algorithm of the certificate.",
			},
			"signature_hash": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Signature hash of the certificate.",
			},
			"extended_key_usage": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Extended key usage of the certificate.",
			},
			"issued_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Issuer of the certificate.",
			},
			"parent_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Parent certificate of this certificate.",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Version of the certificate.",
			},
			"certificate_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the certificate.",
			},
			"days_to_expiry": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of days until the certificate expires (computed by the API).",
			},
			"key_agreement_usage": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the key agreement usage flag is set.",
			},
			"encipher_only_usage": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the encipher only usage flag is set.",
			},
			"key_encipherment_usage": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the key encipherment usage flag is set.",
			},
			"data_encipherment_usage": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the data encipherment usage flag is set.",
			},
			"crl_sign_usage": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the CRL sign usage flag is set.",
			},
			"decipher_only_usage": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the decipher only usage flag is set.",
			},
			"content_commitment_usage": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the content commitment usage flag is set.",
			},
			"digital_signature_usage": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the digital signature usage flag is set.",
			},
			"key_cert_sign_usage": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the key cert sign usage flag is set.",
			},
			"custom_fields": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				Description:      "Any custom fields that will be used in device42.",
				DiffSuppressFunc: suppressCustomFieldsDiffs,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				client := m.(*resty.Client)
				issuedTo := d.Id()

				id, err := getCertificateIdByName(client, issuedTo)
				if err != nil {
					return nil, err
				}

				d.SetId(id)
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceDevice42CertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	issuedTo := d.Get("issued_to").(string)
	formData := map[string]string{}
	formData["issued_to"] = issuedTo

	if v := d.Get("serial_number").(string); v != "" {
		formData["serial_number"] = v
	}
	if v := d.Get("subject").(string); v != "" {
		formData["subject"] = v
	}
	if v := d.Get("san").(string); v != "" {
		formData["san"] = v
	}
	if v := d.Get("valid_from").(string); v != "" {
		formData["valid_from"] = v
	}
	if v := d.Get("valid_to").(string); v != "" {
		formData["valid_to"] = v
	}
	if v := d.Get("vendor").(string); v != "" {
		formData["vendor"] = v
	}
	if v := d.Get("signature_algorithm").(string); v != "" {
		formData["signature_algorithm"] = v
	}
	if v := d.Get("signature_hash").(string); v != "" {
		formData["signature_hash"] = v
	}
	if v := d.Get("extended_key_usage").(string); v != "" {
		formData["extended_key_usage"] = v
	}
	if v := d.Get("issued_by").(string); v != "" {
		formData["issued_by"] = v
	}
	if v := d.Get("parent_cert").(string); v != "" {
		formData["parent_cert"] = v
	}
	if v := d.Get("version").(string); v != "" {
		formData["version"] = v
	}

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/certificates/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42CertificateCreate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42CertificateCreate - Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42CertificateCreate - unable to extract certificate id from response: %s", err)
	}

	d.SetId(strconv.Itoa(id))
	d.Set("certificate_id", id)

	if d.Get("custom_fields") != nil {
		fields := d.Get("custom_fields").(map[string]interface{})
		bulkFields := []string{}

		for k, v := range fields {
			bulkFields = append(bulkFields, fmt.Sprintf("%v:%v", k, v))
		}

		log.Printf("[DEBUG] resourceDevice42CertificateCreate - Pushing custom_fields %s", bulkFields)
		resp, err := client.R().
			SetFormData(map[string]string{
				"id":          strconv.Itoa(id),
				"bulk_fields": strings.Join(bulkFields, ","),
			}).
			SetResult(apiResponse{}).
			Put("/1.0/custom_fields/certificate/")

		if err != nil {
			return err
		}

		r := resp.Result().(*apiResponse)

		if r.Code != 0 {
			return fmt.Errorf("resourceDevice42CertificateCreate - api returned code %d", r.Code)
		}
	}

	return resourceDevice42CertificateRead(d, m)
}

func resourceDevice42CertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	d42Url := fmt.Sprintf("/1.0/certificates/%s/", d.Id())
	log.Printf("[DEBUG] resourceDevice42CertificateRead - reading certificate on target URL: %s", d42Url)

	var result apiCertificateResponse
	apiResp, err := client.R().
		SetResult(&result).
		Get(d42Url)

	if err != nil {
		log.Printf("[WARN] No certificate found: %s", d.Id())
		d.SetId("")
		return nil
	}

	// resty does not surface a 404 as an error, so a deleted/unknown certificate
	// must be detected via the HTTP status code. When it is gone, clear the state.
	if apiResp.StatusCode() == http.StatusNotFound || result.CertificateDetails.ID == 0 {
		log.Printf("[WARN] No certificate found: %s", d.Id())
		d.SetId("")
		return nil
	}

	setCertificateFromDetails(d, result.CertificateDetails)

	return nil
}

func resourceDevice42CertificateUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	issuedTo := d.Get("issued_to").(string)
	formData := map[string]string{}
	formData["issued_to"] = issuedTo

	if d.HasChange("serial_number") {
		formData["serial_number"] = d.Get("serial_number").(string)
	}
	if d.HasChange("subject") {
		formData["subject"] = d.Get("subject").(string)
	}
	if d.HasChange("san") {
		formData["san"] = d.Get("san").(string)
	}
	if d.HasChange("valid_from") {
		formData["valid_from"] = d.Get("valid_from").(string)
	}
	if d.HasChange("valid_to") {
		formData["valid_to"] = d.Get("valid_to").(string)
	}
	if d.HasChange("vendor") {
		formData["vendor"] = d.Get("vendor").(string)
	}
	if d.HasChange("signature_algorithm") {
		formData["signature_algorithm"] = d.Get("signature_algorithm").(string)
	}
	if d.HasChange("signature_hash") {
		formData["signature_hash"] = d.Get("signature_hash").(string)
	}
	if d.HasChange("extended_key_usage") {
		formData["extended_key_usage"] = d.Get("extended_key_usage").(string)
	}
	if d.HasChange("issued_by") {
		formData["issued_by"] = d.Get("issued_by").(string)
	}
	if d.HasChange("parent_cert") {
		formData["parent_cert"] = d.Get("parent_cert").(string)
	}
	if d.HasChange("version") {
		formData["version"] = d.Get("version").(string)
	}

	log.Printf("[DEBUG] resourceDevice42CertificateUpdate - Pushing new informations : %#v", formData)

	resp, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Post("/1.0/certificates/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("resourceDevice42CertificateUpdate - api returned code %d", r.Code)
	}

	log.Printf("[DEBUG] resourceDevice42CertificateUpdate - Result: %#v", r)

	id, err := msgID(r.Msg)
	if err != nil {
		return fmt.Errorf("resourceDevice42CertificateUpdate - unable to extract certificate id from response: %s", err)
	}

	if d.HasChange("custom_fields") {
		updateList := setCustomFields(d)
		for k, v := range updateList {
			log.Printf("[DEBUG] resourceDevice42CertificateUpdate - Pushing custom field %s : %v", k, v)
			resp, err := client.R().
				SetFormData(map[string]string{
					"id":    strconv.Itoa(id),
					"key":   k,
					"value": v.(string),
				}).
				SetResult(apiResponse{}).
				Put("/1.0/custom_fields/certificate/")

			if err != nil {
				return err
			}

			r := resp.Result().(*apiResponse)
			log.Printf("[DEBUG] resourceDevice42CertificateUpdate - Result: %#v", r)
		}
	}

	return resourceDevice42CertificateRead(d, m)
}

func resourceDevice42CertificateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("Deleting certificate %s (ID: %s)", d.Get("issued_to"), d.Id())

	url := fmt.Sprintf("/1.0/certificates/%s/", d.Id())

	resp, err := client.R().
		SetResult(apiDeleteResponse{}).
		Delete(url)

	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return fmt.Errorf("resourceDevice42CertificateDelete - HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	r := resp.Result().(*apiDeleteResponse)
	log.Printf("[DEBUG] resourceDevice42CertificateDelete - Result: %#v", r)
	return nil
}

// setCertificateFromDetails populates the schema fields of the resource/data
// source from an apiCertificateDetails object. Shared by the resource Read and
// the data source Read so the field mapping stays in one place.
func setCertificateFromDetails(d *schema.ResourceData, c apiCertificateDetails) {
	d.Set("issued_to", c.IssuedTo)
	d.Set("certificate_id", c.ID)
	d.Set("serial_number", c.SerialNumber)
	d.Set("subject", c.Subject)
	d.Set("san", c.SAN)
	d.Set("valid_from", c.ValidFrom)
	d.Set("valid_to", c.ValidTo)
	d.Set("vendor", c.Vendor)
	d.Set("signature_algorithm", c.SignatureAlgorithm)
	d.Set("signature_hash", c.SignatureHash)
	d.Set("extended_key_usage", c.ExtendedKeyUsage)
	d.Set("issued_by", c.IssuedBy)
	d.Set("parent_cert", c.ParentCert)
	d.Set("version", string(c.Version))
	d.Set("days_to_expiry", c.DaysToExpiry)
	d.Set("key_agreement_usage", c.KeyAgreementUsage)
	d.Set("encipher_only_usage", c.EncipherOnlyUsage)
	d.Set("key_encipherment_usage", c.KeyEnciphermentUsage)
	d.Set("data_encipherment_usage", c.DataEnciphermentUsage)
	d.Set("crl_sign_usage", c.CRLSignUsage)
	d.Set("decipher_only_usage", c.DecipherOnlyUsage)
	d.Set("content_commitment_usage", c.ContentCommitmentUsage)
	d.Set("digital_signature_usage", c.DigitalSignatureUsage)
	d.Set("key_cert_sign_usage", c.KeyCertSignUsage)
	d.Set("custom_fields", flattenCustomFields(c.CustomFields))
}

// filterCertificateByName returns the single certificate whose issued_to
// exactly matches the given value. Because GET /1.0/certificates/ returns ALL
// certificates, the by-name lookup must be done client-side, erroring when no
// match or more than one match is found.
func filterCertificateByName(certs []apiCertificateDetails, issuedTo string) (apiCertificateDetails, error) {
	var matches []apiCertificateDetails
	for _, c := range certs {
		if c.IssuedTo == issuedTo {
			matches = append(matches, c)
		}
	}

	if len(matches) == 0 {
		return apiCertificateDetails{}, fmt.Errorf("no certificate found with issued_to: %s", issuedTo)
	}

	if len(matches) > 1 {
		return apiCertificateDetails{}, fmt.Errorf("more than one certificate found with issued_to: %s", issuedTo)
	}

	return matches[0], nil
}

func getCertificateIdByName(client *resty.Client, issuedTo string) (string, error) {
	var lookup apiCertificateListResponse
	url := fmt.Sprintf("/1.0/certificates/")
	log.Printf("[DEBUG] Query for certificate name on URL : %s", url)
	_, err := client.R().
		SetResult(&lookup).
		Get(url)

	log.Printf("[DEBUG] Data from certificate lookup : %#v", lookup)
	if err != nil {
		return "", fmt.Errorf("error querying Device42 API: %s", err)
	}

	c, err := filterCertificateByName(lookup.CertificateDetails, issuedTo)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(c.ID, 10), nil
}
