# Table: azure_compute_availability_set

An Availability Set is a logical grouping capability for isolating VM resources from each other when they're deployed.

## Examples

### Availability set basic info

```sql
select
  name,
  platform_fault_domain_count,
  platform_update_domain_count,
  region
from
  azure_compute_availability_set;
```


### List of availability sets which does not use managed disks configuration

```sql
select
  name,
  sku_name
from
  azure_compute_availability_set
where
  sku_name = 'Classic';
```


### List of availability sets without application tag key

```sql
select
  name,
  tags
from
  azure_compute_availability_set
where
  not tags :: JSONB ? 'application';
```
