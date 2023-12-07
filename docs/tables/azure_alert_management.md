---
title: "Steampipe Table: azure_alert_management - Query Azure Alert Management using SQL"
description: "Allows users to query Azure Alert Management, providing a centralized way to monitor and respond to issues across applications and infrastructure."
---

# Table: azure_alert_management - Query Azure Alert Management using SQL

Azure Alert Management is a service within Microsoft Azure that offers a set of tools for monitoring and responding to issues across various Azure resources. It enables users to set up and manage alerts for resources such as virtual machines, databases, web applications, and more. The service helps maintain awareness of the health and performance of Azure resources and facilitates appropriate actions when predefined conditions are met.

## Table Usage Guide

The `azure_alert_management` table provides insights into the alert management system within Microsoft Azure. As a system administrator, you can explore alert-specific details through this table, including alert status, severity, and associated metadata. Use it to identify and respond to potential issues across your Azure resources, ensuring optimal performance and security.

**Important notes:**
- This table offers access to alert management details for the past 30 days. If no value is specified in the query parameter (`time_range`) within the `where` clause, the default value will be set to `1d`(One Day).
- For improved performance, it is advised that you use the optional qual to limit the result set.
- This table supports optional quals. Queries with optional quals are optimized to use Alert Management filters. Optional quals are supported for the following columns:
  - `target_resource`: Filter by the target resource (full ARM ID). The default value selects all resources.
  - `target_resource_type`: Filter by target resource type. The default value selects all resource types.
  - `resource_group`: Filter by target resource group name. The default value selects all resource groups.
  - `alert_rule`: Filter by a specific alert rule. The default value selects all rules.
  - `smart_group_id`: Filter the alerts list by the Smart Group ID. The default value is none.
  - `sort_order`: Sort the query results in ascending or descending order. The default value is 'desc' for time fields and 'asc' for others.
  - `custom_time_range`: Filter by a custom time range in the format (ISO-8601 format). Permissible values are within 30 days from the query time. Either `timeRange` or `customTimeRange` can be used, but not both. The default is none.
  - `sort_by`: Sort the query results by input field. The default value is 'lastModifiedDateTime'. For available fields, refer to the [API documentation](https://learn.microsoft.com/en-us/rest/api/monitor/alertsmanagement/alerts/get-all?tabs=HTTP#alertssortbyfields).
  - `monitor_service`: Filter by the monitor service that generates the alert instance. The default value selects all services. For available services, refer to the [API documentation](https://learn.microsoft.com/en-us/rest/api/monitor/alertsmanagement/alerts/get-all?tabs=HTTP#monitorservice).
  - `monitor_condition`: Filter by the monitor condition, which is either 'Fired' or 'Resolved'. The default value selects all conditions.
  - `severity`: Filter by severity. The default value selects all severities. For details, refer to the [severity documentation](https://learn.microsoft.com/en-us/rest/api/monitor/alertsmanagement/alerts/get-all?tabs=HTTP#severity).
  - `alert_state`: Filter by the state of the alert instance. The default value selects all states. For details, refer to the [alert state documentation](https://learn.microsoft.com/en-us/rest/api/monitor/alertsmanagement/alerts/get-all?tabs=HTTP#alertstate).
  - `time_range`: Filter by the time range, choosing from the listed values in the [API documentation](https://learn.microsoft.com/en-us/rest/api/monitor/alertsmanagement/alerts/get-all?tabs=HTTP#timerange). The default value is 1 day.
v
## Examples

### Basic info
Determine the areas in which Azure's alert management system is currently active. This allows you to understand the overall health and status of your alerts, helping you to manage and respond to potential issues more effectively.

```sql+postgres
select
  name,
  id,
  type,
  target_resource,
  signal_type,
  alert_state,
  monitor_condition
from
  azure_alert_management;
```

```sql+sqlite
select
  name,
  id,
  type,
  target_resource,
  signal_type,
  alert_state,
  monitor_condition
from
  azure_alert_management;
```

### List fired alerts
Explore which alerts have been triggered in your Azure environment to gain insights into potential issues or areas of concern. This helps in proactive problem management and maintaining system stability.

```sql+postgres
select
  name,
  id,
  type,
  signal_type,
  alert_state,
  monitor_service,
  monitor_condition
from
  azure_alert_management
where
  monitor_condition = 'Fired';
```

```sql+sqlite
select
  name,
  id,
  type,
  signal_type,
  alert_state,
  monitor_service,
  monitor_condition
from
  azure_alert_management
where
  monitor_condition = 'Fired';
```

### List alerts within the last 7 days
Explore recent alerts by identifying those that have been generated within the last week. This is useful for maintaining awareness of recent activity and potential issues in your Azure environment.

```sql+postgres
select
  name,
  id,
  target_resource,
  target_resource_type,
  alert_rule,
  time_range
from
  azure_alert_management
where
  time_range = '7d';
```

```sql+sqlite
The given PostgreSQL query does not contain any PostgreSQL-specific functions, data types, or syntax that needs to be converted to SQLite. Therefore, the SQLite query is the same as the PostgreSQL query:

```sql
select
  name,
  id,
  target_resource,
  target_resource_type,
  alert_rule,
  time_range
from
  azure_alert_management
where
  time_range = '7d';
```
```

### List critical alerts
Determine the areas in which critical alerts are present in your Azure resources. This is beneficial for prioritizing and addressing issues that have the highest severity level.

```sql+postgres
select
  name,
  id,
  target_resource,
  target_resource_type,
  severity,
  alert_state,
  monitor_service
from
  azure_alert_management
where
  severity = 'Sev0';
```

```sql+sqlite
select
  name,
  id,
  target_resource,
  target_resource_type,
  severity,
  alert_state,
  monitor_service
from
  azure_alert_management
where
  severity = 'Sev0';
```

### List alerts of VMInsights monitoring service
This example allows users to identify any alerts associated with the VMInsights monitoring service in Azure. This can be useful for administrators who need to quickly assess the status and details of these alerts for troubleshooting or system management purposes.

```sql+postgres
select
  name,
  id,
  target_resource,
  monitor_service,
  alert_rule,
  alert_state,
  source_created_id
from
  azure_alert_management
where
  monitor_service = 'VMInsights';
```

```sql+sqlite
select
  name,
  id,
  target_resource,
  monitor_service,
  alert_rule,
  alert_state,
  source_created_id
from
  azure_alert_management
where
  monitor_service = 'VMInsights';
```