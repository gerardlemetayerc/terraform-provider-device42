package device42

import (
    "context"

    "github.com/device42/go-device42/device42"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDevice42VirtualMachine() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceDevice42VirtualMachineCreate,
        ReadContext:   resourceDevice42VirtualMachineRead,
        UpdateContext: resourceDevice42VirtualMachineUpdate,
        DeleteContext: resourceDevice42VirtualMachineDelete,

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

func resourceDevice42VirtualMachineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*device42.Client)
    name := d.Get("name").(string)
    ipAddress := d.Get("ip_address").(string)
    subnet := d.Get("subnet").(string)
    application := d.Get("application").(string)

    vm, err := client.CreateVirtualMachine(name, ipAddress, subnet, application)
    if err != nil {
        return diag.FromErr(err)
    }

    d.SetId(vm.ID)
    return nil
}

func resourceDevice42VirtualMachineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*device42.Client)
    vm, err := client.GetVirtualMachine(d.Id())
    if err != nil {
        return diag.FromErr(err)
    }

    d.Set("name", vm.Name)
    d.Set("ip_address", vm.IPAddress)
    d.Set("subnet", vm.Subnet)
    d.Set("application", vm.Application)

    return nil
}

func resourceDevice42VirtualMachineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*device42.Client)
    name := d.Get("name").(string)
    ipAddress := d.Get("ip_address").(string)
    subnet := d.Get("subnet").(string)
    application := d.Get("application").(string)

    _, err := client.UpdateVirtualMachine(d.Id(), name, ipAddress, subnet, application)
    if err != nil {
        return diag.FromErr(err)
    }

    return nil
}

func resourceDevice42VirtualMachineDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*device42.Client)

    err := client.DeleteVirtualMachine(d.Id())
    if err != nil {
        return diag.FromErr(err)
    }

    d.SetId("")
    return nil
}
