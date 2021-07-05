# Table: azure_express_route_circuit

ExpressRoute lets you extend your on-premises networks into the Microsoft cloud over a private connection with the help of a connectivity provider. With ExpressRoute, you can establish connections to Microsoft cloud services, such as Microsoft Azure and Microsoft 365.

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

