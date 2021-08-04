# Table: azure_logic_app_workflow

Azure Logic Apps helps you simplify and implement scalable integrations and workflows in the cloud. You can model and automate your process visually as a series of steps known as a workflow in the Logic App Designer.

## Examples

### Basic info

```sql
select
  name,
  id,
  state,
  type
from
  azure_logic_app_workflow;
```

### List disabled workflows

```sql
select
  name,
  id,
  state,
  type
from
  azure_logic_app_workflow
where
  state = 'Disabled';
```

### List suspended workflows

```sql
select
  name,
  id,
  state,
  type
from
  azure_logic_app_workflow
where
  state = 'Suspended';
```
