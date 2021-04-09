# Table: azure_log_alert

Activity log alerts are the alerts that get activated when a new activity log event occurs that matches the conditions specified in the alert.

## Examples

### Basic info

```sql
select
  name,
  jsonb_pretty(activity_log_alert) as activity_log_alert
from
  azure_log_alert;
```

### Ensure that Activity Log Alert exists for Create Policy Assignment

```sql
select
  name,
  id,
  type
from
  azure_log_alert,
  jsonb_array_elements(activity_log_alert -> 'condition' -> 'allOf') as l
where
  l ->> 'equals' = 'Microsoft.Authorization/policyAssignments/write';
```
