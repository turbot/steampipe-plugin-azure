---
title: "Steampipe Table: azure_databricks_workspace - Query Azure Databricks Workspaces using SQL"
description: "Allows users to query Azure Databricks Workspaces, providing insights into the configuration, status, and properties of each workspace."
---

# Table: azure_databricks_workspace - Query Azure Databricks Workspaces using SQL

Azure Databricks Workspace is a feature within Microsoft Azure that offers an interactive workspace for big data analytics and machine learning. It provides a centralized environment for collaborative projects, allowing users to write in multiple languages, visualize data, and share insights. Azure Databricks Workspace supports the full lifecycle of big data analytics, from data preparation to exploration, and from model training to production.

## Table Usage Guide

The `azure_databricks_workspace` table provides insights into Azure Databricks Workspaces within Microsoft Azure. As a data scientist or data analyst, you can explore workspace-specific details through this table, including configuration, status, and properties of each workspace. Use it to uncover information about workspaces, such as their location, SKU, managed private network, and provisioning status.

## Examples

### Basic info
Explore the various Azure Databricks workspaces within your organization to gain insights into their creation dates and associated SKU details. This can be useful for tracking resource usage and understanding your workspace configuration.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that have been recently added to your workspace within the past month. This is especially useful for keeping track of new additions and managing growth in your workspace.

```sql+postgres
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

```sql+sqlite
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
  created_date_time >= datetime('now', '-30 days');
```

### List failed workspaces
Determine the areas in which Azure Databricks workspaces have failed. This can be useful in identifying issues and troubleshooting the workspaces that are not successfully provisioned.

```sql+postgres
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

```sql+sqlite
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
  provisioning_state = 'Failed';
```

### List all encrypted workspaces
Identify instances where workspaces in Azure Databricks are encrypted. This is useful for ensuring data security and compliance with encryption standards.

```sql+postgres
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

```sql+sqlite
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
  json_extract(parameters, '$.Encryption') is not null;
```

### List workspaces that allow public IP
Determine the areas in which Azure Databricks workspaces are configured to allow public IP access. This query can be used to identify potential security vulnerabilities and ensure best practices for data protection.

```sql+postgres
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

```sql+sqlite
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
  json_extract(json_extract(parameters, '$.enableNoPublicIp'), '$.value') = 'false';
```