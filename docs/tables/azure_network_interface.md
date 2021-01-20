# Table: azure_network_interface

A network interface enables an Azure Virtual Machine to communicate with internet, Azure, and on-premises resources.

## Examples

### Basic IP address info

```sql
select
	name,
	ip ->> 'name' as config_name,
	ip -> 'properties' ->> 'privateIPAddress' as private_ip_address,
	ip -> 'properties' ->> 'privateIPAddressVersion' as private_ip_address_version,
	ip -> 'properties' ->> 'privateIPAllocationMethod' as private_ip_address_allocation_method
from
	azure_network_interface
	cross join jsonb_array_elements(ip_configurations) as ip;
```


### Find all network interfaces with private IPs that are in a given subnet (10.66.0.0/16)

```sql
select
	name,
	ip ->> 'name' as config_name,
	ip -> 'properties' ->> 'privateIPAddress' as private_ip_address
from
	azure_network_interface
	cross join jsonb_array_elements(ip_configurations) as ip
where
	ip -> 'properties' ->> 'privateIPAddress' = '10.66.0.0/16';
```


### Security groups attached to each network interface

```sql
select
	name,
	split_part(network_security_group_id, '/', 8) as security_groups
from
	azure_network_interface;
```