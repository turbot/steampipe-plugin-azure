---
title: "Steampipe Table: azure_monitor_activity_log_event - Query Azure Monitor Activity Log Events using SQL"
description: "Allows users to query Activity Log Events in Azure Monitor, providing insights into the operation logs and changes made in Azure resources."
folder: "Monitor"
---

# Table: azure_monitor_activity_log_event - Query Azure Monitor Activity Log Events using SQL

Azure Monitor Activity Log Events is a feature within Microsoft Azure that provides insights into the operational activities within your Azure resources. It enables you to categorize and analyze data about the status, event severity, and operations of your Azure resources. Azure Monitor Activity Log Events helps you stay informed about the activities and operations happening in your Azure environment.

## Table Usage Guide

The `azure_monitor_activity_log_event` table provides insights into the operational activities within Azure Monitor. As a system administrator or a DevOps engineer, explore event-specific details through this table, including event category, event initiation, and associated metadata. Utilize it to uncover information about events, such as those related to service health, resource health, and administrative operations.

**Important notes:**
- This table can provide event details for the previous 90 days.
- For improved performance, it is advised that you use the optional qual `event_timestamp` to limit the result set to a specific time period.
- This table supports optional quals. Queries with optional quals are optimized to use Monitor Activity Log filters. Optional quals are supported for the following columns:
  - `event_timestamp`
  - `resource_group`
  - `correlation_id`
  - `resource_id`
  - `resource_provider_name`

## Examples

### Basic info
Explore the sequence and timing of events in your Azure Monitor Activity Log. This query can be used to gain insights into patterns of activity, identify potential issues, and track changes over time.

```sql+postgres
select
  event_name,
  event_data_id,
  id,
  correlation_id,
  level,
  resource_id,
  event_timestamp
from
  azure_monitor_activity_log_event;
```

```sql+sqlite
select
  event_name,
  event_data_id,
  id,
  correlation_id,
  level,
  resource_id,
  event_timestamp
from
  azure_monitor_activity_log_event;
```

### List events with event-level critical
This example helps identify critical events in your Azure activity log. By doing so, it allows you to promptly respond to potential issues or security threats.

```sql+postgres
select
  event_name,
  id,
  operation_name,
  event_timestamp,
  level,
  caller
from
  azure_monitor_activity_log_event
where
  level = 'EventLevelCritical';
```

```sql+sqlite
select
  event_name,
  id,
  operation_name,
  event_timestamp,
  level,
  caller
from
  azure_monitor_activity_log_event
where
  level = 'EventLevelCritical';
```

### List events that occurred over the last five minutes
Track recent activities in your Azure environment by identifying events that have taken place within the last five minutes. This is useful for real-time monitoring and immediate response to changes or incidents.

```sql+postgres
select
  event_name,
  event_timestamp,
  operation_name,
  resource_id,
  resource_type,
  status
from
  azure_monitor_activity_log_event
where
  event_timestamp >= now() - interval '5 minutes';
```

```sql+sqlite
select
  event_name,
  event_timestamp,
  operation_name,
  resource_id,
  resource_type,
  status
from
  azure_monitor_activity_log_event
where
  event_timestamp >= datetime('now', '-5 minutes');
```

### List ordered events that occurred in the past five to ten minutes
Determine the sequence of events that transpired in the recent past. This can be useful to track and analyze real-time activities, helping to identify patterns or anomalies for prompt action.

```sql+postgres
select
  event_name,
  id,
  submission_timestamp,
  event_timestamp,
  category,
  sub_status
from
  azure_monitor_activity_log_event
where
  event_timestamp between (now() - interval '10 minutes') and (now() - interval '5 minutes')
order by
  event_timestamp asc;
```

```sql+sqlite
select
  event_name,
  id,
  submission_timestamp,
  event_timestamp,
  category,
  sub_status
from
  azure_monitor_activity_log_event
where
  event_timestamp between (datetime('now', '-10 minutes')) and (datetime('now', '-5 minutes'))
order by
  event_timestamp asc;
```

### Get authorization details for events
Determine the authorization details associated with various events to help manage permissions and access control in your Azure environment. This can help in identifying any unauthorized activities or potential security risks.

```sql+postgres
select
  event_name,
  authorization_info ->> 'Action' as authorization_action,
  authorization_info ->> 'Role' as authorization_role,
  authorization_info ->> 'Scope' as authorization_scope
from
  azure_monitor_activity_log_event;
```

```sql+sqlite
select
  event_name,
  json_extract(authorization_info, '$.Action') as authorization_action,
  json_extract(authorization_info, '$.Role') as authorization_role,
  json_extract(authorization_info, '$.Scope') as authorization_scope
from
  azure_monitor_activity_log_event;
```

### Get HTTP request details of events
Explore the specifics of HTTP requests in event logs to identify potential security threats or unusual activity. This could be useful in troubleshooting, security audits, or monitoring network traffic.

```sql+postgres
select
  event_name,
  operation_name,
  event_timestamp,
  http_request ->> 'ClientRequestID' as client_request_id,
  http_request ->> 'ClientIPAddress' as ClientIPAddress,
  http_request ->> 'Method' as method,
  http_request ->> 'URI' as uri
from
  azure_monitor_activity_log_event;
```

```sql+sqlite
select
  event_name,
  operation_name,
  event_timestamp,
  json_extract(http_request, '$.ClientRequestID') as client_request_id,
  json_extract(http_request, '$.ClientIPAddress') as ClientIPAddress,
  json_extract(http_request, '$.Method') as method,
  json_extract(http_request, '$.URI') as uri
from
  azure_monitor_activity_log_event;
```

## Filter examples

### List evens by resource group
Discover the segments that are active within a specific resource group in Azure Monitor's activity log. This can be particularly useful for tracking and managing operations, resources, and statuses associated with specific events.

```sql+postgres
select
  event_name,
  id,
  resource_id,
  operation_name,
  resource_type,
  status
from
  azure_monitor_activity_log_event
where
  resource_group = 'my_rg';
```

```sql+sqlite
select
  event_name,
  id,
  resource_id,
  operation_name,
  resource_type,
  status
from
  azure_monitor_activity_log_event
where
  resource_group = 'my_rg';
```

### List events for a resource provider
Explore the activities associated with a specific resource provider on Azure. This query is useful for tracking operations, event names, and statuses related to a particular network resource provider, helping you understand its activity and performance.

```sql+postgres
select
  event_name,
  id,
  resource_id,
  operation_name,
  resource_provider_name,
  resource_type,
  status
from
  azure_monitor_activity_log_event
where
  resource_provider_name = 'Microsoft.Network';
```

```sql+sqlite
select
  event_name,
  id,
  resource_id,
  operation_name,
  resource_provider_name,
  resource_type,
  status
from
  azure_monitor_activity_log_event
where
  resource_provider_name = 'Microsoft.Network';
```

### List events for a particular resource
Discover the segments that have undergone recent changes in a specific resource within your Azure environment. This is particularly useful for tracking changes and maintaining security compliance.

```sql+postgres
select
  event_name,
  id,
  resource_id,
  event_timestamp,
  correlation_id,
  resource_provider_name
from
  azure_monitor_activity_log_event
where
  resource_id = '/subscriptions/hsjekr16-f95f-4771-bbb5-8237jsa349sl/resourceGroups/my_rg/providers/Microsoft.Network/publicIPAddresses/test-backup-ip';
```

```sql+sqlite
select
  event_name,
  id,
  resource_id,
  event_timestamp,
  correlation_id,
  resource_provider_name
from
  azure_monitor_activity_log_event
where
  resource_id = '/subscriptions/hsjekr16-f95f-4771-bbb5-8237jsa349sl/resourceGroups/my_rg/providers/Microsoft.Network/publicIPAddresses/test-backup-ip';
```