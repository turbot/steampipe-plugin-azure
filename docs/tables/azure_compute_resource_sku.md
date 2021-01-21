# Table: azure_compute_resource_sku

Types of sku available for azure compute resources.

## Examples

### Compute resources sku info

```sql
select
  name,
  tier,
  size,
  family,
  kind
from
  azure_compute_resource_sku;
```


### Azure compute resources and their capacity

```sql
select
  name,
  default_capacity,
  maximum_capacity,
  minimum_capacity
from
  azure_compute_resource_sku;
```


### List of all premium type disks and location

```sql
select
  name,
  resource_type tier,
  l as location
from
  azure_compute_resource_sku,
  jsonb_array_elements_text(locations) as l
where
  resource_type = 'disks'
  and tier = 'Premium';
```