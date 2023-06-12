# Table: azure_dns_zone

Azure DNS zone is used to host the DNS records for a particular domain. Please note that this table only retrieves public DNS zones, use the `azure_private_dns_zone` table for private DNS zones.

## Examples

### Basic info

```sql
select
  name,
  resource_group,
  tags
from
  azure_dns_zone;
```

### List public DNS zones with record sets

```sql
select
  name,
  resource_group
from
  azure_dns_zone
where
  number_of_record_sets > 1;
```

### List public DNS zones with delegated name servers

```sql
select
  name,
  resource_group,
  ns
from
  azure_dns_zone, jsonb_array_elements_text(name_servers) as ns
where
  zone_type = 'Public'
  and ns not like '%.azure-dns.%.';
```
