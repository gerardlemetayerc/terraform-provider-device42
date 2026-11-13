package device42

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// apiSoftwareDetailResponse is the shape of a single "software in use" record
// returned by GET /1.0/software_details/. These are the individual software
// instances detected on devices.
type apiSoftwareDetailResponse struct {
	ID               int64         `json:"id"`
	SoftwareID       int64         `json:"software_id"`
	Software         string        `json:"software"`
	SoftwareAlias    string        `json:"software_alias"`
	DeviceID         int64         `json:"device_id"`
	Device           string        `json:"device"`
	Version          string        `json:"version"`
	FileVersion      string        `json:"file_version"`
	Arch             string        `json:"arch"`
	User             string        `json:"user"`
	UserID           int64         `json:"user_id"`
	Vendor           string        `json:"vendor"`
	VendorID         int64         `json:"vendor_id"`
	Appcomp          string        `json:"appcomp"`
	AppcompID        int64         `json:"appcomp_id"`
	InstallDate      string        `json:"install_date"`
	FirstDetected    string        `json:"first_detected"`
	LastUpdated      string        `json:"last_updated"`
	LicenseUseCount  float64       `json:"license_use_count"`
	CountInLicensing bool          `json:"count_in_licensing"`
	CustomFields     []customField `json:"custom_fields"`
}

// apiSoftwareDetailListResponse is the wrapper returned by
// GET /1.0/software_details/.
type apiSoftwareDetailListResponse struct {
	TotalCount      int64                       `json:"total_count"`
	SoftwareDetails []apiSoftwareDetailResponse `json:"software_details"`
}

func datasourceD42SoftwareDetails() *schema.Resource {
	return &schema.Resource{
		Read:        datasourceD42SoftwareDetailsRead,
		Description: "Read Device42 software-in-use instances, optionally filtered by software or device.",
		Schema: map[string]*schema.Schema{
			"software_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Filter instances to the given software component id.",
			},
			"device_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Filter instances to the given device id.",
			},
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Software-in-use instances matching the filters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"software_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the software component.",
						},
						"software": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the software component.",
						},
						"software_alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alias of the software component.",
						},
						"device_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the device the software is installed on.",
						},
						"device": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the device the software is installed on.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version of the software.",
						},
						"arch": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Architecture of the software.",
						},
						"user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User associated with the instance.",
						},
						"vendor": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vendor of the software.",
						},
						"appcomp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application component associated with the instance.",
						},
						"install_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Installation date of the software.",
						},
						"license_use_count": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "License use count for the instance.",
						},
						"count_in_licensing": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the instance counts towards licensing.",
						},
					},
				},
			},
		},
	}
}

func datasourceD42SoftwareDetailsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	log.Printf("[DEBUG] datasourceD42SoftwareDetailsRead - querying software instances")

	// GET /1.0/software_details/ returns ALL instances in one response, so any
	// filtering is done client-side.
	var resp apiSoftwareDetailListResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get("/1.0/software_details/")
	if err != nil {
		return err
	}
	if apiResp.StatusCode() >= 400 {
		return fmt.Errorf("datasourceD42SoftwareDetailsRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	softwareID := int64(d.Get("software_id").(int))
	deviceID := int64(d.Get("device_id").(int))

	instances := make([]interface{}, 0)
	for _, it := range resp.SoftwareDetails {
		if softwareID > 0 && it.SoftwareID != softwareID {
			continue
		}
		if deviceID > 0 && it.DeviceID != deviceID {
			continue
		}
		instances = append(instances, map[string]interface{}{
			"software_id":        it.SoftwareID,
			"software":           it.Software,
			"software_alias":     it.SoftwareAlias,
			"device_id":          it.DeviceID,
			"device":             it.Device,
			"version":            it.Version,
			"arch":               it.Arch,
			"user":               it.User,
			"vendor":             it.Vendor,
			"appcomp":            it.Appcomp,
			"install_date":       it.InstallDate,
			"license_use_count":  it.LicenseUseCount,
			"count_in_licensing": it.CountInLicensing,
		})
	}
	d.Set("instances", instances)

	// Stable id derived from the filters so distinct queries get distinct ids.
	key := "all"
	if softwareID > 0 {
		key = "sw-" + strconv.FormatInt(softwareID, 10)
	}
	if deviceID > 0 {
		key = key + "-dev-" + strconv.FormatInt(deviceID, 10)
	}
	d.SetId(key)
	return nil
}
