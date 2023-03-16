# Table: azure_private_dns_zone

Azure private DNS zone is used to host the DNS records for a particular domain. Please note that this table only retrieves private DNS zones, use the `azure_dns_zone` table for public DNS zones.

## Examples

### Basic info

```sql
select
  name,
  resource_group,
  tags
from
  azure_private_dns_zone;
```

### List private DNS zones with record sets

```sql
select
  name,
  resource_group
from
  azure_private_dns_zone
where
  number_of_record_sets > 1;
```

### List private DNS zones linked to no virtual networks

```sql
select
  name,
  resource_group
from
  azure_private_dns_zone
where
  number_of_virtual_network_links_with_registration = 0;
```
