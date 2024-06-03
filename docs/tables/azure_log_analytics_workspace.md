---
title: "Steampipe Table: azure_log_analytics_workspace - Query Azure Log Analytics Workspaces using SQL"
description: "Allows users to query Azure Log Analytics Workspaces, providing insights into the configuration and properties of Log Analytics Workspaces within their Azure environment."
---

# Table: azure_log_analytics_workspace - Query Azure Log Analytics Workspaces using SQL

Azure Log Analytics Workspaces are environments for managing log data collected from various sources. These workspaces are essential for centralizing and analyzing log data to monitor and gain insights into the operational health and performance of Azure resources.

## Table Usage Guide

The `azure_log_analytics_workspace` table provides insights into the properties and configurations of Log Analytics Workspaces within your Azure environment. System administrators can explore details such as the workspace's SKU, retention settings, network access configurations, and various features. Utilize this table to manage and optimize your log data collection and analysis processes.

## Examples

### Basic Info

Retrieve basic information about your Log Analytics Workspaces, including their names, IDs, and locations.

```sql
select
  name,
  id,
  location
from
  azure_log_analytics_workspace;
```

### List Workspaces with Retention Period Greater than 30 Days

Identify workspaces where the log retention period exceeds 30 days.

```sql
select
  name,
  id,
  retention_in_days
from
  azure_log_analytics_workspace
where
  retention_in_days > 30;
```

### Get Workspaces with Specific Features Enabled

Find workspaces that have data export enabled.

```sql
select
  name,
  id,
  enable_data_export
from
  azure_log_analytics_workspace
where
  enable_data_export = true;
```

### Get Workspaces with Data Export Enabled

Find workspaces that have data export enabled.

```sql
select
  name,
  id,
  enable_data_export
from
  azure_log_analytics_workspace
where
  enable_data_export = true;
```

### Identify Workspaces with Disabled Local Auth

List workspaces where non-AAD based authentication is disabled.

```sql
select
  name,
  id,
  disable_local_auth
from
  azure_log_analytics_workspace
where
  disable_local_auth = true;
```

### Workspaces with Private Link Scoped Resources

Retrieve workspaces that have linked private link scope resources.

```sql
select
  name,
  id,
  private_link_scoped_resources
from
  azure_log_analytics_workspace
where
  private_link_scoped_resources is not null;
```

### Workspaces with Force CMK for Query Enabled

Find workspaces where customer-managed keys are mandatory for query management.

```sql
select
  name,
  id,
  force_cmk_for_query
from
  azure_log_analytics_workspace
where
  force_cmk_for_query = true;
```
