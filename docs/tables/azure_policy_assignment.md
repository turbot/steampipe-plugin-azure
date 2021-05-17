# Table: azure_policy_assignment

Policy assignments are used by Azure Policy to define which resources are assigned which policies or initiatives in a subscription.

## Examples

### Basic info

```sql
select
  id,
  policy_definition_id,
  name,
  type
from
  azure_policy_assignment;
```

### Get SQL auditing and threat detection monitoring status for the subscription

```sql
select
  id,
  policy_definition_id,
  display_name,
  parameters -> 'sqlAuditingMonitoringEffect' -> 'value' as sqlAuditingMonitoringEffect
from
  azure_policy_assignment;
```

### Get SQL encryption monitoring status for the subscription

```sql
select
  id,
  policy_definition_id,
  display_name,
  parameters -> 'sqlEncryptionMonitoringEffect' -> 'value' as sqlEncryptionMonitoringEffect
from
  azure_policy_assignment;
```
