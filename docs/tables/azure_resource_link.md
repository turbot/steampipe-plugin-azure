# Table: azure_resource_link

Linking is a feature of the Resource Manager. It enables declaring relationships between resources even if they do not reside in the same resource group.

## Examples

### Basic Info

```sql
select
  name,
  id,
  type,
  source_id,
  target_id
from
  azure_resource_link;
```

### List resource links with virtual machines

```sql
select
  name,
  id,
  source_id,
  target_id
from
  azure_resource_link
where
  source_id LIKE '%virtualmachines%';
```
