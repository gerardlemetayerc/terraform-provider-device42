package device42

import (
	"errors"
	"fmt"
)

func (c *Client) CreateDevice(name string, ip string, subnet string, application string) error {
	vm, err := c.createVirtualMachine(name, ip, subnet, application)
	if err != nil {
		return err
	}

	// TODO: add code to create device in Device42

	fmt.Printf("Created device %s with IP %s\n", name, vm.IP)
	return nil
}

func (c *Client) UpdateDevice(name string, ip string, subnet string, application string) error {
	vm, err := c.updateVirtualMachine(name, ip, subnet, application)
	if err != nil {
		return err
	}

	// TODO: add code to update device in Device42

	fmt.Printf("Updated device %s with IP %s\n", name, vm.IP)
	return nil
}

func (c *Client) DeleteDevice(name string) error {
	vm, err := c.deleteVirtualMachine(name)
	if err != nil {
		return err
	}

	// TODO: add code to delete device in Device42

	fmt.Printf("Deleted device %s with IP %s\n", name, vm.IP)
	return nil
}

func (c *Client) createVirtualMachine(name string, ip string, subnet string, application string) (*VirtualMachine, error) {
	if name == "" || ip == "" || subnet == "" || application == "" {
		return nil, errors.New("missing required parameter(s)")
	}

	// TODO: add code to create virtual machine in Device42

	vm := &VirtualMachine{
		Name:        name,
		IP:          ip,
		Subnet:      subnet,
		Application: application,
	}
	return vm, nil
}

func (c *Client) updateVirtualMachine(name string, ip string, subnet string, application string) (*VirtualMachine, error) {
	if name == "" || ip == "" || subnet == "" || application == "" {
		return nil, errors.New("missing required parameter(s)")
	}

	// TODO: add code to update virtual machine in Device42

	vm := &VirtualMachine{
		Name:        name,
		IP:          ip,
		Subnet:      subnet,
		Application: application,
	}
	return vm, nil
}

func (c *Client) deleteVirtualMachine(name string) (*VirtualMachine, error) {
	if name == "" {
		return nil, errors.New("missing required parameter(s)")
	}

	// TODO: add code to delete virtual machine in Device42

	vm := &VirtualMachine{
		Name: name,
	}
	return vm, nil
}
