# Table: azure_logic_app_workflow

A workflow is a series of steps that defines a task or process. Each workflow starts with a single trigger, after which you must add one or more actions.

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
