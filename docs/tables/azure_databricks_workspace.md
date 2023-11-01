# Table: azure_databricks_workspace

A workspace is an environment for accessing all of your Azure Databricks assets. A workspace organizes objects (notebooks, libraries, dashboards, and experiments) into folders and provides access to data objects and computational resources.

## Examples

### Basic info

```sql
select
  name,
  id,
  workspace_id,
  workspace_url,
  created_date_time,
  sku
from
  azure_databricks_workspace;
```

### List workspaces created in the last 30 days

```sql
select
  name,
  id,
  workspace_id,
  workspace_url,
  created_date_time,
  sku
from
  azure_databricks_workspace
where
  created_date_time >= now() - interval '30' day;
```

### List failed workspaces

```sql
select
  name,
  id,
  workspace_id,
  workspace_url,
  created_date_time,
  sku
from
  azure_databricks_workspace;
where
  provisioning_state = 'Failed';
```

### List all encrypted workspaces

```sql
select
  name,
  id,
  workspace_id,
  workspace_url,
  created_date_time,
  sku
from
  azure_databricks_workspace
where
  parameters -> 'Encryption' is not null;
```

### List workspaces that allow public IP

```sql
select
  name,
  id,
  workspace_id,
  workspace_url,
  created_date_time,
  sku
from
  azure_databricks_workspace
where
  parameters -> 'enableNoPublicIp' ->> 'value' = 'false';
```
