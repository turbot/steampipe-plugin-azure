# Table: azure_automation_variable

Variable assets are values that are available to all runbooks and DSC configurations in your Automation account. You can manage them from the Azure portal, from PowerShell, within a runbook, or in a DSC configuration.

## Examples

### Basic info

```sql
select
  id,
  name,
  account_name,
  type,
  is_encrypted,
  value
from
  azure_automation_variable;
```

### List variables that are unencrypted

```sql
select
  id,
  name,
  account_name,
  type,
  is_encrypted,
  value
from
  azure_automation_variable
where
  not is_encrypted;
```

### List variables created in last 30 days

```sql
select
  id,
  name,
  account_name,
  creation_time,
  type,
  is_encrypted,
  value
from
  azure_automation_variable
where
  creation_time >= now() - interval '30' day;
```