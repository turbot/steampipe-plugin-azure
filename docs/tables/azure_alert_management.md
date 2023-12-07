---
title: "Steampipe Table: azure_alert_management - Query Azure Alert Management using SQL"
description: "Allows users to query Azure Alert Management, providing a centralized way to monitor and respond to issues across applications and infrastructure."
---

# Table: azure_alert_management - Query Azure Alert Management using SQL

Azure Alert Management is a service within Microsoft Azure that offers a set of tools for monitoring and responding to issues across various Azure resources. It enables users to set up and manage alerts for resources such as virtual machines, databases, web applications, and more. The service helps maintain awareness of the health and performance of Azure resources and facilitates appropriate actions when predefined conditions are met.

## Table Usage Guide

The `azure_alert_management` table provides insights into the alert management system within Microsoft Azure. As a system administrator, you can explore alert-specific details through this table, including alert status, severity, and associated metadata. Use it to identify and respond to potential issues across your Azure resources, ensuring optimal performance and security.

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