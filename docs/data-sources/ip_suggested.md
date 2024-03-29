---
page_title: "device42_device Data Source - terraform-provider-device42"
subcategory: "network management"
description: |-
  
---

# device42_ip_suggested (Data Source)

This datasource permit to query free IP from a given subnet. Each call will query an IP.
Reservation need to be done using ip datasource.



## Exemple 

```hcl

datasource "device42_ip_suggested" "searchFreeIP" {
  subnet_id    = 1
}

resource "device42_ip" "ip" {
  subnet 	          = device42_subnet.myNewNetwork.name
  ip        	      = data.device42_suggestedIp.ip.ip
  available 	      = "no"
  vrf_group_id 	    = data.device42_subnet.subnet.vrf_group_id
  lifecycle {
    ignore_changes = [
      ip, subnet
    ]
  }
 depends_on = [device42_subnet.myNewNetwork]
}
```


## Schema

### Optional

- `subnet_name` (String) The name of the subnet.
- `subnet_id`(Int) The network ID

### Computed

- `ip` (String) Free IP returned by Device42 API.