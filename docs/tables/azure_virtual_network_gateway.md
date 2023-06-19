# Table: azure_virtual_network_gateway

A virtual network gateway is used to establish secure, cross-premises connectivity.

## Examples

### Basic info

```sql
select
  name,
  id,
  enable_bgp,
  region,
  resource_group
from
  azure_virtual_network_gateway;
```

### List network gateways with no connections

```sql
select
  name,
  id,
  enable_bgp,
  region,
  resource_group
from
  azure_virtual_network_gateway
where
   gateway_connections is null;
```
