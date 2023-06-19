# Table: azure_route_table

Azure Route Tables allows to create network routes so that CloudGen Firewall VM can handle the traffic both between the subnets and to the Internet.

## Examples

### List of subnets associated with route table

```sql
select
  name,
  split_part(subnet ->> 'id', '/', 11) subnet,
  region
from
  azure_route_table
  cross join jsonb_array_elements(subnets) as subnet;
```

### List of route tables where route propagation is enabled

```sql
select
  name,
  disable_bgp_route_propagation,
  region
from
  azure_route_table
where
  not disable_bgp_route_propagation;
```

### Route info of each routes table

```sql
select
  name,
  route ->> 'name' route_name,
  route -> 'properties' ->> 'addressPrefix' address_prefix,
  route -> 'properties' ->> 'nextHopType' next_hop_type
from
  azure_route_table
  cross join jsonb_array_elements(routes) as route;
```
