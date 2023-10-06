# Table: azure_alert_management

Azure Alert Management is a service and set of tools within Microsoft Azure that allows you to monitor and respond to issues across your applications and infrastructure. It provides a centralized way to set up and manage alerts for various Azure resources, including virtual machines, databases, web applications, and more. Azure Alert Management helps you stay informed about the health and performance of your Azure resources and take appropriate actions when predefined conditions are met.

**Important notes:**

- This table offers access to alert management details for the past 30 days. If no value is specified in the query parameter (`time_range`) within the WHERE clause, the default value will be set to `1d`(One Day).
- For improved performance, it is advised that you use the optional qual to limit the result set.
- This table supports optional quals. Queries with optional quals are optimised to use Alert Management filters. Optional quals are supported for the following columns:

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

## Examples

### Basic info

```sql
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

```sql
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

### List alerts with in last 7 days

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

### List critical alerts

```sql
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

```sql
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