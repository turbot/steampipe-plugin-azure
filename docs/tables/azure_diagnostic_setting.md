# Table: azure_diagnostic_setting

Azure diagnostic settings used to send platform logs and metrics to different destinations.

## Examples

### Basic info

```sql
select
  name,
  id,
  type
from
  azure_diagnostic_setting;
```

### Ensure Diagnostic Setting captures Alert category

```sql
select
  l ->> 'category' as category,
  l ->> 'enabled' as enabled
from
  azure_diagnostic_setting,
  jsonb_array_elements(diagnostic_settings -> 'logs') as l
where
  diagnostic_settings is not null
  and l ->> 'category' = 'Alert'
  and l ->> 'enabled' = 'true';
```
