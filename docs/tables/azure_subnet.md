# Table: azure_subnet

A subnet is a range of IP addresses in the VNet. You can divide a VNet into multiple subnets for organization and security

## Examples

### Virtual network and IP address range of each subnet

```sql
select
	name,
	virtual_network_name,
	address_prefix,
	location,
	resource_group
from
	azure_subnet;
```


### Route table associated with each subnet

```sql
select
	st.name subnet_name,
	st.virtual_network_name,
	rt.name route_table_name,
	jsonb_array_elements(rt.routes) -> 'properties' ->> 'addressPrefix' as route_address_prefix,
	jsonb_array_elements(rt.routes) -> 'properties' ->> 'nextHopType' as route_next_hop_type
from
	azure_route_table as rt
	join azure_subnet st on rt.id = st.route_table_id;
```


### Network security group associated with each subnet

```sql
select
	name subnet_name,
	virtual_network_name,
	split_part(network_security_group_id, '/', 9) as network_security_name
from
	azure_subnet;
```


### Service endpoints info of each subnet

```sql
select
	name,
	endpoint -> 'locations' as location,
	endpoint -> 'service' as service
from
	azure_subnet
	cross join jsonb_array_elements(service_endpoints) as endpoint;
```
