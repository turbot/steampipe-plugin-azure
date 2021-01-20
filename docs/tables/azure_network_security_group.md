# Table: azure_network_security_group

A network security group contains security rules that allow or deny inbound network traffic to, or outbound network traffic from, several types of Azure resources.

## Examples

### Subnets and network interfaces attached to the network security groups

```sql
select
	name,
	split_part(nic ->> 'id', '/', 9) network_interface,
	split_part(vn ->> 'id', '/', 9) virtual_network,
	split_part(vn ->> 'id', '/', 11) subnets
from
	azure_network_security_group
	cross join jsonb_array_elements(network_interfaces) as nic,
	jsonb_array_elements(subnets) as vn;
```


### List the network security groups whose inbound is not restricted from the internet

```sql
select
	name,
	sg ->> 'name' as sg_name,
	sg -> 'properties' ->> 'access' as access,
	sg -> 'properties' ->> 'description' as description,
	sg -> 'properties' ->> 'destinationPortRange' as destination_port_range,
	sg -> 'properties' ->> 'direction' as direction,
	sg -> 'properties' ->> 'priority' as priority,
	sg -> 'properties' ->> 'sourcePortRange' as source_port_range,
	sg -> 'properties' ->> 'protocol' as protocol
from
	azure_network_security_group
	cross join jsonb_array_elements(security_rules) as sg
where
	(
		sg -> 'properties' ->> 'sourcePortRange' = '*'
		and sg -> 'properties' ->> 'destinationPortRange' = '*'
    and sg -> 'properties' ->> 'access' = 'Allow'
	);
```


### Default security group rules info

```sql
select
	name,
	sg -> 'name' as sg_name,
	sg -> 'properties' ->> 'access' as access,
	sg -> 'properties' ->> 'description' as description,
	sg -> 'properties' ->> 'destinationPortRange' as destination_port_range,
	sg -> 'properties' ->> 'direction' as direction,
	sg -> 'properties' ->> 'priority' as priority,
	sg -> 'properties' ->> 'sourcePortRange' as source_port_range,
	sg -> 'properties' ->> 'protocol' as protocol
from
	azure_network_security_group
	cross join jsonb_array_elements(default_security_rules) as sg;
```