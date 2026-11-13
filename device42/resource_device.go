package device42

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type customField struct {
	Key   string      `json:"key"`
	Notes interface{} `json:"notes"`
	Value interface{} `json:"value"`
}

type apiDeviceReadResponse struct {
	AssetNo         string        `json:"asset_no"`
	Building        string        `json:"building"`
	Category        string        `json:"category"`
	CustomFields    []customField `json:"custom_fields"`
	Customer        string        `json:"customer"`
	DeviceID        int64         `json:"device_id"`
	ID              int64         `json:"id"`
	InService       bool          `json:"in_service"`
	IsItBladeHost   string        `json:"is_it_blade_host"`
	IsItSwitch      string        `json:"is_it_switch"`
	IsItVirtualHost string        `json:"is_it_virtual_host"`
	Manufacturer    string        `json:"manufacturer"`
	Name            string        `json:"name"`
	Notes           string        `json:"notes"`
	OS              string        `json:"os"`
	SerialNo        string        `json:"serial_no"`
	ServiceLevel    string        `json:"service_level"`
	Tags            []interface{} `json:"tags"`
	Type            string        `json:"type"`
}

// apiCode tolerates Device42's inconsistent `code` encoding: some endpoints
// return a JSON number (0), others a quoted string ("0"). A plain int64 fails
// to unmarshal the string form, which breaks the operation.
type apiCode int64

func (c *apiCode) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		v, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse code %q: %w", s, err)
		}
		*c = apiCode(v)
		return nil
	}
	var n int64
	if err := json.Unmarshal(b, &n); err == nil {
		*c = apiCode(n)
		return nil
	}
	return fmt.Errorf("cannot unmarshal code value: %s", string(b))
}

type apiResponse struct {
	Code apiCode       `json:"code"`
	Msg  []interface{} `json:"msg"`
}

type apiArchiveResponse struct {
	Code apiCode       `json:"code"`
	Msg  []interface{} `json:"msg"`
}

// UnmarshalJSON tolerates Device42's inconsistent archive response shapes.
// The /2.0/devices/{id}/archive/ endpoint returns a bare JSON array on some
// instances (e.g. ["Object archived",11,"name",true,true]) instead of the
// usual {"code":0,"msg":[...]} object. A bare array is treated as Msg.
func (r *apiArchiveResponse) UnmarshalJSON(b []byte) error {
	trimmed := bytes.TrimSpace(b)
	if len(trimmed) > 0 && trimmed[0] == '[' {
		return json.Unmarshal(trimmed, &r.Msg)
	}
	type alias apiArchiveResponse
	var aux alias
	if err := json.Unmarshal(trimmed, &aux); err != nil {
		return err
	}
	*r = apiArchiveResponse(aux)
	return nil
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
			"building": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the building that contains this device.",
			},
			"building_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the building that contains this device.",
			},
			"storage_room": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the storage room that contains this device.",
			},
			"storage_room_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the storage room that contains this device.",
			},
			"datacenter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the data center that contains this device.",
			},
			"asset_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Asset number of the device.",
			},
			"serial_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Serial number of the device.",
			},
			"manufacturer": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Hardware manufacturer of the device.",
			},
			"hardware": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Hardware model name of the device.",
			},
			"hardware_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Hardware model ID of the device.",
			},
			"os_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Operating system name of the device.",
			},
			"in_service": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Whether the device is in service (yes/no).",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Notes for the device.",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Tags of the device (comma separated).",
			},
			"tags_remove": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Tags to remove from the device (comma separated).",
			},
			"objectcategory": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Object category name of the device.",
			},
			"objectcategory_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Object category ID of the device.",
			},
			"customer": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Customer associated with the device.",
			},
			"customer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Customer ID associated with the device.",
			},
			"virtual_host_device": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Virtual host device name for this device.",
			},
			"network_device": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Whether the device is a network switch (yes/no).",
			},
			"blade_chassis": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Whether the device is a blade chassis (yes/no).",
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
	formData := map[string]string{
		"name":          name,
		"type":          deviceType,
		"service_level": serviceLevel,
	}
	if v := d.Get("building").(string); v != "" {
		formData["building"] = v
	}
	if v, ok := d.GetOk("building_id"); ok {
		formData["building_id"] = strconv.Itoa(v.(int))
	}
	if v := d.Get("storage_room").(string); v != "" {
		formData["storage_room"] = v
	}
	if v, ok := d.GetOk("storage_room_id"); ok {
		formData["storage_room_id"] = strconv.Itoa(v.(int))
	}
	if v := d.Get("datacenter").(string); v != "" {
		formData["datacenter"] = v
	}
	if v := d.Get("asset_no").(string); v != "" {
		formData["asset_no"] = v
	}
	if v := d.Get("serial_no").(string); v != "" {
		formData["serial_no"] = v
	}
	if v := d.Get("manufacturer").(string); v != "" {
		formData["manufacturer"] = v
	}
	if v := d.Get("hardware").(string); v != "" {
		formData["hardware"] = v
	}
	if v := d.Get("hardware_id").(string); v != "" {
		formData["hardware_id"] = v
	}
	if v := d.Get("os_name").(string); v != "" {
		formData["os_name"] = v
	}
	if v := d.Get("in_service").(string); v != "" {
		formData["in_service"] = v
	}
	if v := d.Get("notes").(string); v != "" {
		formData["notes"] = v
	}
	if v := d.Get("tags").(string); v != "" {
		formData["tags"] = v
	}
	if v := d.Get("tags_remove").(string); v != "" {
		formData["tags_remove"] = v
	}
	if v := d.Get("objectcategory").(string); v != "" {
		formData["objectcategory"] = v
	}
	if v := d.Get("objectcategory_id").(string); v != "" {
		formData["objectcategory_id"] = v
	}
	if v := d.Get("customer").(string); v != "" {
		formData["customer"] = v
	}
	if v := d.Get("customer_id").(string); v != "" {
		formData["customer_id"] = v
	}
	if v := d.Get("virtual_host_device").(string); v != "" {
		formData["virtual_host_device"] = v
	}
	if v := d.Get("network_device").(string); v != "" {
		formData["network_device"] = v
	}
	if v := d.Get("blade_chassis").(string); v != "" {
		formData["blade_chassis"] = v
	}

	resp, err := client.R().
		SetFormData(formData).
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

	var resp apiDeviceReadResponse
	apiResp, err := client.R().
		SetResult(&resp).
		Get(d42Url)

	if err != nil {
		log.Printf("[WARN] No device found: %s", d.Id())
		d.SetId("")
		return err
	}
	if apiResp.StatusCode() == 404 {
		log.Printf("[WARN] No device found: %s", d.Id())
		d.SetId("")
		return nil
	}
	if apiResp.StatusCode() >= 400 {
		return fmt.Errorf("resourceDevice42DeviceRead - HTTP %d: %s", apiResp.StatusCode(), apiResp.String())
	}

	// Update main fields
	d.Set("name", resp.Name)
	d.Set("type", resp.Type)
	d.Set("service_level", resp.ServiceLevel)
	d.Set("building", resp.Building)
	d.Set("asset_no", resp.AssetNo)
	d.Set("serial_no", resp.SerialNo)
	d.Set("manufacturer", resp.Manufacturer)
	d.Set("customer", resp.Customer)
	d.Set("notes", resp.Notes)
	d.Set("os_name", resp.OS)
	d.Set("tags", joinBuildingList(resp.Tags))
	d.Set("objectcategory", resp.Category)
	d.Set("network_device", resp.IsItSwitch)
	d.Set("blade_chassis", resp.IsItBladeHost)
	d.Set("in_service", yesNo(resp.InService))

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
	if d.HasChange("building") {
		formData["building"] = d.Get("building").(string)
	}
	if d.HasChange("building_id") {
		formData["building_id"] = strconv.Itoa(d.Get("building_id").(int))
	}
	if d.HasChange("storage_room") {
		formData["storage_room"] = d.Get("storage_room").(string)
	}
	if d.HasChange("storage_room_id") {
		formData["storage_room_id"] = strconv.Itoa(d.Get("storage_room_id").(int))
	}
	if d.HasChange("datacenter") {
		formData["datacenter"] = d.Get("datacenter").(string)
	}
	if d.HasChange("asset_no") {
		formData["asset_no"] = d.Get("asset_no").(string)
	}
	if d.HasChange("serial_no") {
		formData["serial_no"] = d.Get("serial_no").(string)
	}
	if d.HasChange("manufacturer") {
		formData["manufacturer"] = d.Get("manufacturer").(string)
	}
	if d.HasChange("hardware") {
		formData["hardware"] = d.Get("hardware").(string)
	}
	if d.HasChange("hardware_id") {
		formData["hardware_id"] = d.Get("hardware_id").(string)
	}
	if d.HasChange("os_name") {
		formData["os_name"] = d.Get("os_name").(string)
	}
	if d.HasChange("in_service") {
		formData["in_service"] = d.Get("in_service").(string)
	}
	if d.HasChange("notes") {
		formData["notes"] = d.Get("notes").(string)
	}
	if d.HasChange("tags") {
		formData["tags"] = d.Get("tags").(string)
	}
	if d.HasChange("tags_remove") {
		formData["tags_remove"] = d.Get("tags_remove").(string)
	}
	if d.HasChange("objectcategory") {
		formData["objectcategory"] = d.Get("objectcategory").(string)
	}
	if d.HasChange("objectcategory_id") {
		formData["objectcategory_id"] = d.Get("objectcategory_id").(string)
	}
	if d.HasChange("customer") {
		formData["customer"] = d.Get("customer").(string)
	}
	if d.HasChange("customer_id") {
		formData["customer_id"] = d.Get("customer_id").(string)
	}
	if d.HasChange("virtual_host_device") {
		formData["virtual_host_device"] = d.Get("virtual_host_device").(string)
	}
	if d.HasChange("network_device") {
		formData["network_device"] = d.Get("network_device").(string)
	}
	if d.HasChange("blade_chassis") {
		formData["blade_chassis"] = d.Get("blade_chassis").(string)
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
		if err != nil {
			log.Printf("[ERROR] Failed to archive Device ID %s", d.Id())
			return err
		}
		if resp1.StatusCode() >= 400 && resp1.StatusCode() != 404 {
			return fmt.Errorf("resourceDevice42DeviceDelete - archive HTTP %d: %s", resp1.StatusCode(), resp1.String())
		}
		r1 := resp1.Result().(*apiArchiveResponse)
		log.Printf("[DEBUG] Result: %#v", r1)
		return nil
	} else {
		resp2, err := client.R().
			SetResult(apiResponse{}).
			Delete(url)
		if err != nil {
			return err
		}
		if resp2.StatusCode() >= 400 && resp2.StatusCode() != 404 {
			return fmt.Errorf("resourceDevice42DeviceDelete - HTTP %d: %s", resp2.StatusCode(), resp2.String())
		}
		r2 := resp2.Result().(*apiResponse)
		log.Printf("[DEBUG] Result: %#v", r2)
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

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func getDeviceIdFromHostname(client *resty.Client, hostname string) (string, error) {
	var apiDeviceLookUp struct {
		TotalCount int `json:"total_count"`
		Devices    []struct {
			DeviceID int64  `json:"device_id"`
			Name     string `json:"name"`
		} `json:"devices"`
	}
	url := fmt.Sprintf("/2.0/devices/?name=%s", url.QueryEscape(hostname))
	log.Printf("[DEBUG] Query for hostname on URL : %s", url)
	_, err := client.R().
		SetResult(&apiDeviceLookUp).
		Get(url)

	log.Printf("[DEBUG] Data from device lookup : %#v", apiDeviceLookUp)
	if err != nil {
		return "", fmt.Errorf("error querying Device42 API: %s", err)
	}

	if len(apiDeviceLookUp.Devices) == 0 {
		return "", fmt.Errorf("no device found with hostname: %s", hostname)
	}

	return strconv.FormatInt(apiDeviceLookUp.Devices[0].DeviceID, 10), nil
}
