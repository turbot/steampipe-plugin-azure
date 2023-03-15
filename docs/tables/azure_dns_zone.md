# Table: azure_dns_zone

Azure DNS zone is used to host the DNS records for a particular domain.

## Examples

### Basic info

```sql
select
  name,
  id,
  zone_type
from
  azure_dns_zone;
```

### List private DNS zones

```sql
select
  name,
  id,
  tags
from
  azure_dns_zone
where
  zone_type = 'Private';
```

### List public DNS zones with delegated name servers

```sql
select
  name,
  id,
  ns
from
  azure_dns_zone, jsonb_array_elements_text(name_servers) as ns
where
  zone_type = 'Public'
  and ns not like '%.azure-dns.%.';
```
