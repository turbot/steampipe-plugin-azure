# Table: azure_nat_gateway

A network interface enables an Azure Virtual Machine to communicate with internet, Azure, and on-premises resources.

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

### Public IP address info

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


### List subnet details associated with nat gatways

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