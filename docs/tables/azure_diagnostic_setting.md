# Table: azure_diagnostic_setting

Azure diagnostic settings are used to send platform logs and metrics to different destinations.

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

### List diagnostic settings that capture Alert category logs

```sql
select
  name,
  id,
  type
from
  azure_diagnostic_setting,
  jsonb_array_elements(logs) as l
where
  l ->> 'category' = 'Alert'
  and l ->> 'enabled' = 'true';
```

### List diagnostic settings that capture Security category logs

```sql
select
  name,
  id,
  type
from
  azure_diagnostic_setting,
  jsonb_array_elements(logs) as l
where
  l ->> 'category' = 'Security'
  and l ->> 'enabled' = 'true';
```

### List diagnostic settings that capture Policy category logs

```sql
select
  name,
  id,
  type
from
  azure_diagnostic_setting,
  jsonb_array_elements(logs) as l
where
  l ->> 'category' = 'Policy'
  and l ->> 'enabled' = 'true';
```

### List diagnostic settings that capture Administrative category logs

```sql
select
  name,
  id,
  type
from
  azure_diagnostic_setting,
  jsonb_array_elements(logs) as l
where
  l ->> 'category' = 'Administrative'
  and l ->> 'enabled' = 'true';
```
