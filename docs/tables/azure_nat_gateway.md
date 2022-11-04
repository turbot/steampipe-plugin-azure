# Table: azure_nat_gateway

NAT gateway provides outbound internet connectivity for one or more subnets of a virtual network. Once NAT gateway is associated to a subnet, NAT provides source network address translation (SNAT) for that subnet. NAT gateway specifies which static IP addresses virtual machines use when creating outbound flows.

## Examples

### Basic info

```sql
select
  name,
  id,
  provisioning_state,
  sku_name,
  type
from
  azure_nat_gateway;
```

### List public IP address details for each nat gateway

```sql
select
  n.name,
  i.ip_address as ip_address,
  i.ip_configuration_id as ip_configuration_id,
  i.public_ip_address_version as public_ip_address_version,
  i.public_ip_allocation_method as public_ip_allocation_method
from
  azure_nat_gateway as n,
  azure_public_ip as i,
  jsonb_array_elements(n.public_ip_addresses) as ip
where
  ip ->> 'id' = i.id;
```


### List subnet details associated with each nat gateway

```sql
select
  n.name as name,
  s.name as subnet_name,
  s.virtual_network_name as virtual_network_name
from
  azure_nat_gateway as n,
  azure_subnet as s,
  jsonb_array_elements(n.subnets) as sb
where
  sb ->> 'id' = s.id;
```
