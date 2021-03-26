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
  jsonb_pretty(activity_log_alert) as activity_log_alert
from
  azure_log_alert
where
  jsonb_path_exists(
    activity_log_alert,
    '$.** ? (@.type() == "string" && @ like_regex "Microsoft.Authorization/policyAssignments/write")'
  )
```