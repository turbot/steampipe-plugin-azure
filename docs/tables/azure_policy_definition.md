# Table: azure_policy_definition

Azure Policy establishes conventions for resources. Policy definitions describe resource compliance conditions and the effect to take if a condition is met. A condition compares a resource property field or a value to a required value.

## Examples

### Basic info

```sql
select
  id,
  name,
  display_name,
  type,
  jsonb_pretty(policy_rule) as policy_rule
from
  azure_policy_definition
```

### Get the policy definition by display name

```sql
select
  id,
  name,
  display_name,
  type,
  jsonb_pretty(policy_rule) as policy_rule
from
  azure_policy_definition
where
  display_name = 'Private endpoint connections on Batch accounts should be enabled';
```
