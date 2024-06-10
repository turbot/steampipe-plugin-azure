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

Retrieve basic information about your Log Analytics Workspaces, including their names, IDs, and locations. This helps in keeping track of all available workspaces within your Azure environment.

```sql+postgres
select
  name,
  id,
  location
from
  azure_log_analytics_workspace;
```

```sql+sqlite
select
  name,
  id,
  location
from
  azure_log_analytics_workspace;
```

### List workspaces with retention period greater than 30 days

Identify workspaces where the log retention period exceeds 30 days. This can be useful for compliance and data retention policy enforcement.

```sql+postgres
select
  name,
  id,
  retention_in_days
from
  azure_log_analytics_workspace
where
  retention_in_days > 30;
```

```sql+sqlite
select
  name,
  id,
  retention_in_days
from
  azure_log_analytics_workspace
where
  retention_in_days > 30;
```

### Get workspaces that have data export enabled

Find workspaces that have data export enabled. This is essential for monitoring data export activities and ensuring that important data is being transferred as expected.

```sql+postgres
select
  name,
  id,
  enable_data_export
from
  azure_log_analytics_workspace
where
  enable_data_export = true;
```

```sql+sqlite
select
  name,
  id,
  enable_data_export
from
  azure_log_analytics_workspace
where
  enable_data_export = true;
```

### Identify workspaces with local auth disabled

List workspaces where non-AAD based authentication is disabled. This information is crucial for maintaining secure access controls and adhering to organizational security policies.

```sql+postgres
select
  name,
  id,
  disable_local_auth
from
  azure_log_analytics_workspace
where
  disable_local_auth = true;
```

```sql+sqlite
select
  name,
  id,
  disable_local_auth
from
  azure_log_analytics_workspace
where
  disable_local_auth = true;
```

### Workspaces with private link scoped resources

Retrieve workspaces that have linked private link scope resources. This helps in understanding the private network configurations and ensuring secure communication within your Azure environment.

```sql+postgres
select
  name,
  id,
  private_link_scoped_resources
from
  azure_log_analytics_workspace
where
  private_link_scoped_resources is not null;
```

```sql+sqlite
select
  name,
  id,
  private_link_scoped_resources
from
  azure_log_analytics_workspace
where
  private_link_scoped_resources is not null;
```

### Workspaces with force CMK for query enabled

Find workspaces where customer-managed keys are mandatory for query management. This is important for organizations that require additional security measures for data encryption and query operations.

```sql+postgres
select
  name,
  id,
  force_cmk_for_query
from
  azure_log_analytics_workspace
where
  force_cmk_for_query = true;
```

```sql+sqlite
select
  name,
  id,
  force_cmk_for_query
from
  azure_log_analytics_workspace
where
  force_cmk_for_query = true;
```
