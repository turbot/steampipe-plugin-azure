# Table: azure_express_route_circuit

An ExpressRoute circuit represents a logical connection between your on-premises infrastructure and Microsoft cloud services through a connectivity provider. You can order multiple ExpressRoute circuits. Each circuit can be in the same or different regions, and can be connected to your premises through different connectivity providers.

## Examples

### Basic info

```sql
select
  name,
  id,
  allow_classic_operations,
  circuit_provisioning_state
from
  azure_express_route_circuit;
```


### List express route circuits which have global reach enabled

```sql
select
  name,
  sku_tier,
  sku_name
from
  azure_express_route_circuit
where
  global_reach_enabled;
```


### List the premium category express route circuits

```sql
select
  name,
  sku_tier,
  sku_name
from
  azure_express_route_circuit
where
  sku_tier = 'Premium';
```

