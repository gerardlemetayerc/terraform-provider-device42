package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type PostResponseMsg struct {
	Code int       `json:"code"`
	Msg  [5]string `json:"msg"`
}

type Device struct {
	Last_updated  string        `json:"last_updated"`
	Orientation   int           `json:"orientation"`
	Ip_addresses  []interface{} `json:"ip_addresses"`
	Serial_no     string        `json:"serial_no"`
	Id            string        `json:"id"`
	Service_level string        `json:"service_level"`
	Uuid          string        `json:"uuid"`
	Name          string        `json:"name"`
	Mac_addresses []interface{} `json:"mac_addresses"`
	Os            string        `json:"os"`
}

func resourceDevice() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeviceCreate,
		Read:   resourceDeviceRead,
		Update: resourceDeviceUpdate,
		Delete: resourceDeviceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func readDeviceFromId(_device_id string) Device {

	_url := HostEnv + "/api/1.0/devices/id/" + _device_id
	client := &http.Client{}

	req, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", AuthString)
	resp, err := client.Do(req)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var responseDevice Device
	json.Unmarshal(bodyBytes, &responseDevice)

	return responseDevice

}

func resourceDeviceCreate(d *schema.ResourceData, m interface{}) error {
	const (
		_url = HostEnv + "/api/1.0/device"
	)

	client := &http.Client{}

	req, err := http.NewRequest("POST", _url, nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", AuthString)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var responseObject PostResponseMsg

	json.Unmarshal(bodyBytes, &responseObject)
	if responseObject.Code == 0 {
		device := readDeviceFromId(responseObject.Msg[1])
		d.Set("name", device.Name)
		d.Set("id", device.Id)
	}
	return nil
}

func resourceDeviceRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDeviceUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDeviceDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
