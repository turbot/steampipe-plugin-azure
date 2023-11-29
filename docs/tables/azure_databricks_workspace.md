---
title: "Steampipe Table: azure_databricks_workspace - Query Azure Databricks Workspaces using SQL"
description: "Allows users to query Azure Databricks Workspaces."
---

# Table: azure_databricks_workspace - Query Azure Databricks Workspaces using SQL

Azure Databricks is an Apache Spark-based analytics platform optimized for the Microsoft Azure cloud services platform. It provides a collaborative environment for data scientists, data engineers, and business analysts to work together. Azure Databricks allows you to build, train, and deploy AI solutions at scale.

## Table Usage Guide

The 'azure_databricks_workspace' table provides insights into Databricks Workspaces within Azure Databricks. As a data scientist or engineer, explore workspace-specific details through this table, including configurations, locations, and associated metadata. Utilize it to uncover information about workspaces, such as those with specific configurations, the relationships between workspaces, and the verification of workspace settings. The schema presents a range of attributes of the Databricks Workspace for your analysis, like the workspace ID, creation date, provisioning state, and associated tags.

## Examples

### Basic info
Explore the basic information about Azure Databricks workspaces, such as their names and IDs. This can be useful to understand the distribution and usage of workspaces across your Azure environment.

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
Explore which workspaces have been established within the past month. This is useful for keeping track of recent additions and understanding the growth of your workspace environment.

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
Discover the segments that have experienced unsuccessful provisioning in Azure Databricks to understand where issues might have occurred. This is useful in identifying potential problems in your setup that may need troubleshooting.

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
Discover the segments that utilize encrypted workspaces in Azure Databricks. This is beneficial in assessing the security measures in place within your organization's data processing environment.

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
Identify the Azure Databricks workspaces that are configured to allow public IP access. This can be useful for assessing potential security risks and ensuring compliance with company policies.

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