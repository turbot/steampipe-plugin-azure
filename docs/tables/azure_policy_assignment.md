# Table: azure_policy_assignment

Azure policy assignment retrieves the information of all policy assignments associated with the given subscription.

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

### Policy assignment status for sqlAuditingMonitoringEffect parameter

```sql
select
  id,
  policy_definition_id,
  display_name,
  parameters -> 'sqlAuditingMonitoringEffect' -> 'value' as sqlAuditingMonitoringEffect
from
  azure_policy_assignment;
```

### Policy assignment status for sqlEncryptionMonitoringEffect parameter

```sql
select
  id,
  policy_definition_id,
  display_name,
  parameters -> 'sqlEncryptionMonitoringEffect' -> 'value' as sqlEncryptionMonitoringEffect
from
  azure_policy_assignment;
```
