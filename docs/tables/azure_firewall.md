# Table: azure_firewall

Azure Firewall is a managed, cloud-based network security service that protects your Azure Virtual Network resources. It's a fully stateful firewall as a service with built-in high availability and unrestricted cloud scalability.

## Examples

### Azure firewall location and availability zone count info

```sql
select
  name,
  location,
  jsonb_array_length(availability_zones) availability_zones_count
from
  azure_firewall;
```


### Basic IP configuration info

```sql
select
  name,
  ip #> '{properties, privateIPAddress}' private_ip_address,
  ip #> '{properties, privateIPAllocationMethod}' private_ip_allocation_method,
  split_part(
    ip -> 'properties' -> 'publicIPAddress' ->> 'id',
    '/',
    9
  ) public_ip_address_id,
  split_part(ip -> 'properties' ->> 'subnet', '/', 9) virtual_network
from
  azure_firewall
  cross join jsonb_array_elements(ip_configurations) as ip;
```


### List the premium category firewalls

```sql
select
  name,
  sku_tier,
  sku_name
from
  azure_firewall
where
  sku_tier = 'Premium';
```


### List of firewalls where threat intel mode is off

```sql
select
  name,
  threat_intel_mode
from
  azure_firewall
where
  threat_intel_mode = 'Off';
```
