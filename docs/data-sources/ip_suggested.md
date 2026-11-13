---
page_title: "device42_ip_suggested Data Source - terraform-provider-device42"
subcategory: "Network Management"
description: |-
  Deprecated: use the device42_ip resource instead, which computes the next available IP itself.
---

# device42_ip_suggested (Data Source)

~> **Deprecated** This data source is deprecated and will be removed in a future release. Its functionality has been moved **directly into the `device42_ip` resource**, which now computes the next available IP itself.

## Why it is deprecated

The original design separated the "free IP suggestion" from the "IP reservation":

- `device42_ip_suggested` (this data source) only returned a **suggested** free IP (read-only, `GET /1.0/suggest_ip/...`).
- The caller then had to create a `device42_ip` resource to actually **reserve** it.

This split was error-prone:

- **No reservation.** This data source does **not** reserve the IP it returns. The IP stays free until a separate `device42_ip` resource is created.
- **Unstable plans.** Because nothing is reserved by this data source, running `terraform plan` again can return a different (or the same re-suggested) free IP, leading to churn between plans.
- **Two-step coupling.** Users had to write two blocks and correlate them manually (see the example below).

The suggestion logic has therefore been **folded into the `device42_ip` resource**, which now reserves the IP it computes in a single Create. Use `device42_ip` and you no longer need this data source.

## Example Usage (deprecated)

```hcl
# Deprecated: prefer computing the IP directly in device42_ip.
data "device42_ip_suggested" "searchFreeIP" {
  subnet_id = 1
}

resource "device42_ip" "ip" {
  subnet       = device42_subnet.myNewNetwork.name
  ip           = data.device42_ip_suggested.searchFreeIP.ip
  available    = "no"
  vrf_group_id = data.device42_subnet.subnet.vrf_group_id

  lifecycle {
    ignore_changes = [
      ip, subnet
    ]
  }

  depends_on = [device42_subnet.myNewNetwork]
}
```

## Argument Reference (deprecated)

### Optional

- `subnet_name` (String) The name of the subnet.
- `subnet_id` (Int) The network ID.

### Computed

- `id` (String) IP as an ID.
- `ip` (String) Calculated free IP.

## Related

- [`device42_ip`](../resources/ip.md) — the resource that now computes and reserves the next available IP.
- [`device42_subnet_suggested`](subnet_suggested.md) — the subnet counterpart (suggests and reserves a sub-subnet).