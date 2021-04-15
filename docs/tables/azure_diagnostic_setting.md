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

### List diagnostic settings that captures Alert category logs

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

### List diagnostic settings that captures Security category logs

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

### List diagnostic settings that captures Policy category logs

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

### List diagnostic settings that captures Administrative category logs

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
