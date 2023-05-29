# Table: azure_automation_account

 Automation accounts allow you to isolate your Automation resources, runbooks, assets, and configurations from the resources of other accounts. You can use Automation accounts to separate resources into separate logical environments or delegated responsibilities.

## Examples

### Basic info

```sql
select
  name,
  id,
  resource_group,
  type
from
  azure_automation_account;
```

### List accounts that are created in last 30 days

```sql
select
  name,
  id,
  resource_group,
  type,
  creation_time
from
  azure_automation_account
where
  creation_time >= now() - interval '30' day;
```

### List accounts that are suspended

```sql
select
  name,
  id,
  resource_group,
  type,
  creation_time,
  state
from
  azure_automation_account
where
  state = 'AccountStateSuspended';
```
