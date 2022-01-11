---
layout: "websupportsk"
page_title: "Websupportsk: websupportsk_dns_record"
sidebar_current: "docs-websupportsk-resource-dns-record"
description: |-
  Manages a DNS record in the Websupportsk DNS Service
---

# websupportsk\_dns\_record

Manages a DNS record set in the Websupportsk DNS Service.

## Example Usage

### Create a new subdomain DNS record

```hcl
resource "websupportsk_dns_record" "example_record" {
  zone_name  = "example.com"
  ttl        = 6000
  type       = "A"
  name       = "subdomain"
  content    = "10.0.0.0"
}
```

## Argument Reference

The following arguments are supported:

* `zone_name` - (Required) The name of the zone in which to create the record.
  Changing this creates a new DNS record.

* `type` - (Required) The type of record set. Accepts only: "A", "AAAA", "MX", "ANAME", "CNAME", "NS", "TXT", "SRV".
  Changing this creates a new DNS record.

* `name` - (Required) The name of the record.

* `content` - (Required) Content of the DNS record. Example: for record of type "A" it must be a valid IP address.

* `ttl` - (Optional) The time to live (TTL) of the record.

* `prio` - (Required FOR type("MX", "SRV") | else IGNORED) Priority value of the record.

* `port` - (Required FOR type("SRV") | else IGNORED) Port value of the record.

* `weight` - (Required FOR type("SRV") | else IGNORED) Weight value of the record.

* `note` - (Optional) Text note for the record.



## Attributes Reference

The following attributes are exported:

* `zone_name` - See Argument Reference above.
* `type` - See Argument Reference above.
* `name` - See Argument Reference above.
* `content` - See Argument Reference above.
* `ttl` - See Argument Reference above.
* `prio` - See Argument Reference above.
* `weight` - See Argument Reference above.
* `note` - See Argument Reference above.

## Import

Records can be imported using a composite ID formed of zone name and record ID, e.g.

```
$ terraform import websupportsk_dns_record.example_record websupport.sk/123
```

where:

* `websupport.sk` - The name of the zone
* `123` - record ID as returned by [API](https://rest.websupport.sk/docs/v1.zone#records)
