# Table: azure_security_center_sub_assessment

Azure security center sub assessment details for the subscription.

## Examples

### Basic info

```sql
select
  id,
  name,
  display_name,
  type,
  category
from
  azure_security_center_sub_assessment;
```

### List unhealthy sub assessment details

```sql
select
  name,
  type,
  category,
  status
from
  azure_security_center_sub_assessment
where
  lower(status) like '%unhealthy%';
```
