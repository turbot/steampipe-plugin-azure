# Table: azure_virtual_network

A virtual network is a network where all devices, servers, virtual machines, and data centers that are connected are done so through software and wireless technology.

## Examples

### List of virtual networks where DDoS(Distributed Denial of Service attacks) Protection is not enabled

```sql
select
	name,
	enable_ddos_protection,
	region,
	resource_group
from
	azure_virtual_network
where
	not enable_ddos_protection;
```


### CIDR list for each virtual network

```sql
select
	name,
	jsonb_array_elements_text(address_prefixes) as address_block
from
	azure_virtual_network;
```


### List VPCs with public CIDR blocks

```sql
select
	name,
	cidr_block,
	region,
	resource_group
from
	azure_virtual_network
	cross join jsonb_array_elements_text(address_prefixes) as cidr_block
where
	not cidr_block :: cidr < <= '10.0.0.0/16'
	and not cidr_block :: cidr < <= '192.168.0.0/16'
	and not cidr_block :: cidr < <= '172.16.0.0/12';
```


### Subnet details associated with the virtual network

```sql
select
	name,
	subnet ->> 'name' as subnet_name,
	subnet -> 'properties' ->> 'addressPrefix' as address_prefix,
	subnet -> 'properties' ->> 'privateEndpointNetworkPolicies' as private_endpoint_network_policies,
	subnet -> 'properties' ->> 'privateLinkServiceNetworkPolicies' as private_link_service_network_policies,
	subnet -> 'properties' ->> 'serviceEndpoints' as service_endpoints,
	split_part(subnet -> 'properties' ->> 'routeTable', '/', 9) as route_table
from
	azure_virtual_network
	cross join jsonb_array_elements(subnets) as subnet;
```
