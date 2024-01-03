package device42

import (
	"fmt"
	"os"
	"log"
	"strconv"
	"strings"
	"net/url"
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type customField struct {
	Key   string      `json:"key"`
	Notes interface{} `json:"notes"`
	Value interface{} `json:"value"`
}

type apiDeviceReadResponse struct {
	CustomFields            []customField `json:"custom_fields"`
	DeviceID                int64         `json:"device_id"`
	ID                      int64         `json:"id"`
	Name                    string        `json:"name"`
	ServiceLevel            string        `json:"service_level"`
	Tags                    []interface{} `json:"tags"`
	Type                    string        `json:"type"`
}

type apiResponse struct {
	Code int64         `json:"code"`
	Msg  []interface{} `json:"msg"`
}

type apiArchiveResponse struct {
	Code int64 `json:"code"`
}

func resourceD42Device() *schema.Resource {
	return &schema.Resource{
		Create: resourceDevice42DeviceCreate,
		Read:   resourceDevice42DeviceRead,
		Update: resourceDevice42DeviceUpdate,
		Delete: resourceDevice42DeviceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The hostname of the device.",
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "The type of the device. " +
					"Valid values are 'physical', 'virtual', 'unknown', 'cluster' (default is virtual)",
			},
			"service_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "d42null",
				Description: "Service Level of the device (default is d42null).",
			},
			"custom_fields": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				Description:      "Any custom fields that will be used in device42.",
				DiffSuppressFunc: suppressCustomFieldsDiffs,
			},
			"archive_on_destroy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Archive device on destroy action",
			},
		},

		Importer: &schema.ResourceImporter{
		    StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		        client := m.(*resty.Client)
		        hostname := d.Id()
		
		        deviceId, err := getDeviceIdFromHostname(client, hostname)
		        if err != nil {
		            return nil, err
		        }
		
		        d.SetId(deviceId)
		        return []*schema.ResourceData{d}, nil
		    },
		},
	}
}

func resourceDevice42DeviceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	name := d.Get("name").(string)
	deviceType := d.Get("type").(string)
	serviceLevel := d.Get("service_level").(string)

	if deviceType == "" {
		deviceType = "virtual"
	}
	resp, err := client.R().
		SetFormData(map[string]string{
			"name":          name,
			"type":          deviceType,
			"service_level": serviceLevel,
		}).
		SetResult(apiResponse{}).
		Post("/2.0/devices/")

	if err != nil {
		return err
	}

	r := resp.Result().(*apiResponse)

	if r.Code != 0 {
		return fmt.Errorf("API returned code %d", r.Code)
	}

	log.Printf("[DEBUG] Result: %#v", r)

	if len(r.Msg) < 1 {
		str := fmt.Sprintf("%v", r.Msg)
		return fmt.Errorf("please check account permission or credentials : returned msg %s", str)
	}

	id := int(r.Msg[1].(float64))

	if d.Get("custom_fields") != nil {
		fields := d.Get("custom_fields").(map[string]interface{})
		bulkFields := []string{}

		for k, v := range fields {
			bulkFields = append(bulkFields, fmt.Sprintf("%v:%v", k, v))
		}

		resp, err := client.R().
			SetFormData(map[string]string{
				"name":        name,
				"bulk_fields": strings.Join(bulkFields, ","),
			}).
			SetResult(apiResponse{}).
			Put("/1.0/device/custom_field/")

		if err != nil {
			return err
		}

		r := resp.Result().(*apiResponse)

		if r.Code != 0 {
			return fmt.Errorf("API returned code %d", r.Code)
		}
	}

	// Only set ID after all conditions are successfull
	d.SetId(strconv.Itoa(id))

	return resourceDevice42DeviceRead(d, m)
}

// Permit to read data about device in Device42. It use deviceID to query data.
func resourceDevice42DeviceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	d42Url := fmt.Sprintf("/1.0/devices/id/%s/", d.Id())
	log.Printf("[DEBUG] reading Device on target URL: %s", d42Url)
	tfLog := os.Getenv("TF_LOG")
	if tfLog == "DEBUG" {
		client.SetDebug(true)
	}

	var resp apiDeviceReadResponse
	_, err := client.R().
		SetResult(&resp).
		Get(d42Url)

	if err != nil {
		log.Printf("[DEBUG] HTTP Response Status Code: %d", resp.StatusCode())
		log.Printf("[WARN] No device found: %s", d.Id())
		d.SetId("")
		return err
	}

	// Update main fields
	d.Set("name", resp.Name)
	d.Set("type", resp.Type)
	d.Set("service_level", resp.ServiceLevel)

	// Add custom fields
	fields := flattenCustomFields(resp.CustomFields)
	d.Set("custom_fields", fields)

	return nil
}

func resourceDevice42DeviceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	formData := map[string]string{}
	name := d.Get("name").(string)
	d42url := fmt.Sprintf("/2.0/devices/%s/", d.Id())

	if d.HasChange("type") {
		formData["type"] = d.Get("type").(string)
	}

	if d.HasChange("service_level") {
		formData["service_level"] = d.Get("service_level").(string)
	}
	log.Printf("[DEBUG] resourceDevice42DeviceUpdate - Updating : %s", d42url)
	log.Printf("[DEBUG] resourceDevice42DeviceUpdate - Pushing new informations : %#v", formData)

	_, err := client.R().
		SetFormData(formData).
		SetResult(apiResponse{}).
		Put(d42url)
	
	if err != nil {
		return err
	}

	if d.HasChange("custom_fields") {
		updateList := setCustomFields(d)
		for k, v := range updateList {
			resp, err := client.R().
				SetFormData(map[string]string{
					"name":  name,
					"key":   k,
					"value": v.(string),
				}).
				SetResult(apiResponse{}).
				Put("/1.0/device/custom_field/")

			if err != nil {
				return err
			}

			r := resp.Result().(*apiResponse)
			log.Printf("[DEBUG] Result: %#v", r)
		}
	}
	return resourceDevice42DeviceRead(d, m)
}

func resourceDevice42DeviceDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*resty.Client)
	archiveOnDetroy := false
	if d.Get("archive_on_destroy").(string) != "" {
		archiveOnDetroy, _ = strconv.ParseBool(d.Get("archive_on_destroy").(string))
	}
	log.Printf("Deleting device %s (UUID: %s)", d.Get("name"), d.Id())

	url := fmt.Sprintf("/1.0/devices/%s/", d.Id())

	if archiveOnDetroy {
		url := fmt.Sprintf("/2.0/devices/%s/archive/", d.Id())
		resp1, err := client.R().
			SetResult(apiArchiveResponse{}).
			Post(url)
		r1 := resp1.Result().(*apiArchiveResponse)
		log.Printf("[DEBUG] Result: %#v", r1)
		if err != nil {
			log.Printf("[ERROR] Failed to archive Device ID %s", d.Id())
			return err
		}
		return nil
	} else {
		resp2, err := client.R().
			SetResult(apiResponse{}).
			Delete(url)
		r2 := resp2.Result().(*apiResponse)
		log.Printf("[DEBUG] Result: %#v", r2)
		if err != nil {
			return err
		}
		return nil
	}
}

func setCustomFields(d *schema.ResourceData) map[string]interface{} {
	updatedFields := make(map[string]interface{})
	if d.HasChange("custom_fields") {
		oldRaw, newRaw := d.GetChange("custom_fields")
		old := oldRaw.(map[string]interface{})
		new := newRaw.(map[string]interface{})
		for k, v := range new {
			if old[k] != v {
				log.Printf("[DEBUG] Change to custom field: %s, Old Value: '%s', New Value: '%s'", k, old[k], v)
				updatedFields[k] = v
			}
		}
	}
	return updatedFields
}

func flattenCustomFields(in []customField) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for _, x := range in {
		out[x.Key] = x.Value
	}
	return out
}

func suppressCustomFieldsDiffs(k, old, new string, d *schema.ResourceData) bool {
	field := strings.TrimPrefix(k, "custom_fields.")
	setFields := d.Get("custom_fields").(map[string]interface{})
	if _, ok := setFields[field]; ok {
		return false
	}
	return true
}

func getDeviceIdFromHostname(client *resty.Client, hostname string) (string, error) {
    var apiResp struct {
        TotalCount int `json:"total_count"`
        Devices    []struct {
            DeviceID int64  `json:"device_id"`
            Name     string `json:"name"`
        } `json:"devices"`
    }

    _, err := client.R().
        SetResult(&apiResp).
        Get(fmt.Sprintf("/api/2.0/devices/?name=%s", url.QueryEscape(hostname)))

    if err != nil {
        return "", fmt.Errorf("error querying Device42 API: %s", err)
    }

    if len(apiResp.Devices) == 0 {
        return "", fmt.Errorf("no device found with hostname: %s", hostname)
    }

    return strconv.FormatInt(apiResp.Devices[0].DeviceID, 10), nil
}
