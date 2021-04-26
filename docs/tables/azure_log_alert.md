# Table: azure_log_alert

Activity log alerts are the alerts that get activated when a new activity log event occurs that matches the conditions specified in the alert.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  enabled
from
  azure_log_alert;
```

### List log alerts that check for create policy assignment events

```sql
select
  name,
  id,
  type
from
  azure_log_alert,
  jsonb_array_elements(condition -> 'allOf') as l
where
  l ->> 'equals' = 'Microsoft.Authorization/policyAssignments/write';
```
