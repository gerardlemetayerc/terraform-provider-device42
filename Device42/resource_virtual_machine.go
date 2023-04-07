package main

import (
	"context"
	"fmt"

	"github.com/device42/device42-go/device42"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualMachineCreate,
		Read:   resourceVirtualMachineRead,
		Update: resourceVirtualMachineUpdate,
		Delete: resourceVirtualMachineDelete,

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

func resourceVirtualMachineCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*device42.Client)
	name := d.Get("name").(string)
	ipAddress := d.Get("ip_address").(string)
	subnet := d.Get("subnet").(string)
	application := d.Get("application").(string)

	vm := device42.VirtualMachine{
		Name:         name,
		IP:           ipAddress,
		Subnet:       subnet,
		BusinessApp:  application,
	}

	newVM, _, err := client.CreateVirtualMachine(context.Background(), vm)
	if err != nil {
		return fmt.Errorf("error creating virtual machine: %s", err)
	}

	d.SetId(fmt.Sprintf("%d", newVM.ID))

	return resourceVirtualMachineRead(d, m)
}

func resourceVirtualMachineRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*device42.Client)

	id := d.Id()

	vm, _, err := client.GetVirtualMachineByID(context.Background(), id)
	if err != nil {
		return fmt.Errorf("error retrieving virtual machine: %s", err)
	}

	d.Set("name", vm.Name)
	d.Set("ip_address", vm.IP)
	d.Set("subnet", vm.Subnet)
	d.Set("application", vm.BusinessApp)

	return nil
}

func resourceVirtualMachineUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*device42.Client)

	id := d.Id()

	vm := device42.VirtualMachine{
		ID:           id,
		Name:         d.Get("name").(string),
		IP:           d.Get("ip_address").(string),
		Subnet:       d.Get("subnet").(string),
		BusinessApp:  d.Get("application").(string),
	}

	_, _, err := client.UpdateVirtualMachine(context.Background(), vm)
	if err != nil {
		return fmt.Errorf("error updating virtual machine: %s", err)
	}

	return resourceVirtualMachineRead(d, m)
}

func resourceVirtualMachineDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*device42.Client)

	id := d.Id()

	err := client.DeleteVirtualMachine(context.Background(), id)
	if err != nil {
		return fmt.Errorf("error deleting virtual machine: %s", err)
	}

	d.SetId("")

	return nil
}
