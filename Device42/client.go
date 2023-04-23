package device42

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Client struct {
	BaseURL    string
	HttpClient *http.Client
}

type VirtualMachine struct {
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Subnet  string `json:"subnet"`
	AppName string `json:"app_name"`
}

func (c *Client) CreateVirtualMachine(vm *VirtualMachine) error {
	url := c.BaseURL + "/device"

	jsonBody, err := json.Marshal(vm)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("CreateVirtualMachine failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) ReadVirtualMachine(id int) (*VirtualMachine, error) {
	url := c.BaseURL + "/device/" + strconv.Itoa(id)

	resp, err := c.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ReadVirtualMachine failed with status code: %d", resp.StatusCode)
	}

	var vm VirtualMachine
	err = json.NewDecoder(resp.Body).Decode(&vm)
	if err != nil {
		return nil, err
	}

	return &vm, nil
}

func (c *Client) UpdateVirtualMachine(id int, vm *VirtualMachine) error {
	url := c.BaseURL + "/device/" + strconv.Itoa(id)

	jsonBody, err := json.Marshal(vm)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("UpdateVirtualMachine failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) DeleteVirtualMachine(id int) error {
	url := c.BaseURL + "/device/" + strconv.Itoa(id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DeleteVirtualMachine failed with status code: %d", resp.StatusCode)
	}

	return nil
}
