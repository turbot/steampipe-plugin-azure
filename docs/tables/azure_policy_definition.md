# Table: azure_policy_definition

Azure Policy establishes conventions for resources. Policy definitions describe resource compliance conditions and the effect to take if a condition is met. A condition compares a resource property field or a value to a required value.

## Examples

### Basic info

```sql
select
  id,
  name,
  display_name,
  type
from
  azure_policy_definition
where
  display_name = 'Private endpoint connections on Batch accounts should be enabled';
```

### Get the policy definitions for a particular subscription

```sql
select
  id,
  name,
  display_name,
  type
from
  azure_policy_definition
where
  subscription_id = '3510ae4d-530b-497d-8f30-53b9616fc6c1';
```

