# Table: azure_machine_learning_workspace

The workspace is the top-level resource for Azure Machine Learning, providing a centralized place to work with all the artifacts you create when you use Azure Machine Learning. The workspace keeps a history of all training runs, including logs, metrics, output, and a snapshot of your scripts. You use this information to determine which training run produces the best model.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  provisioning_state
from
  azure_machine_learning_workspace;
```

### List system assigned identity type workspace

```sql
select
  name,
  id,
  type,
  identity ->> 'type' as identity_type
from
  azure_machine_learning_workspace
where
  identity ->> 'type' = 'SystemAssigned';
```

### List key vault used by workspaces with soft deletion disabled

```sql
select
  m.name as workspace_name,
  m.id as workspace_id,
  v.soft_delete_enabled
from
  azure_machine_learning_workspace as m,
  azure_key_vault as v
where
  lower(m.key_vault) = lower(v.id) and not v.soft_delete_enabled;
```