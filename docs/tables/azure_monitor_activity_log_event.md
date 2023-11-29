---
title: "Steampipe Table: azure_monitor_activity_log_event - Query Azure Monitor Activity Log Events using SQL"
description: "Allows users to query Azure Monitor Activity Log Events"
---

# Table: azure_monitor_activity_log_event - Query Azure Monitor Activity Log Events using SQL

Azure Monitor collects, analyzes, and acts on telemetry data from your Azure and on-premises environments. It helps you understand how your applications are performing and proactively identifies issues affecting them and the resources they depend on. Activity Log Events in Azure Monitor provides insight into subscription-level events that have occurred in Azure.

## Table Usage Guide

The 'azure_monitor_activity_log_event' table provides insights into activity log events within Azure Monitor. As a DevOps engineer, explore event-specific details through this table, including event categories, event data, and associated metadata. Utilize it to uncover information about events, such as those related to service health, resource management, and security. The schema presents a range of attributes of the activity log event for your analysis, like the event timestamp, resource group, event ID, and associated tags.

## Examples

### Basic info
Explore the Azure Monitor activity log to gain insights into the events occurring in your Azure resources. This query can help you understand the scope and impact of each event, making it easier to manage your resources and respond to issues.

```sql
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
Identify instances where critical events have occurred in your Azure Monitor activity log. This could be useful in troubleshooting and understanding the severity of issues within your Azure environment.

```sql
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
This query is used to monitor recent activities within a system, specifically events that have occurred in the last five minutes. It's useful for real-time tracking and immediate response to any critical changes or anomalies in the system.

```sql
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

### List ordered events that occurred in the past five to ten minutes
Explore the sequence of events that happened in the recent past to understand any system changes or unusual activity. This allows for real-time monitoring and swift response to any unexpected events.

```sql
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

### Get authorization details for events
Determine the authorization details associated with specific events to gain insights into the actions, roles, and scopes involved. This can be beneficial for understanding the security context of activities within your Azure environment.

```sql
select
  event_name,
  authorization_info ->> 'Action' as authorization_action,
  authorization_info ->> 'Role' as authorization_role,
  authorization_info ->> 'Scope' as authorization_scope
from
  azure_monitor_activity_log_event;
```

### Get HTTP request details of events
Analyze the details of HTTP requests associated with specific events to understand their operational patterns and time-stamps. This can help in tracking the client's request ID, IP address, and the methods used, which could be beneficial in enhancing security and monitoring network traffic.

```sql
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

## Filter examples

### List evens by resource group
Explore the activities within a specific resource group in Azure Monitor, helping you understand the operations and status of resources for effective management and troubleshooting.

```sql
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
Determine the areas in which specific events are occurring for a particular resource provider in Azure. This can help in analyzing the operation status and type of resources being used, which can be useful for optimizing resource allocation and troubleshooting issues.

```sql
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
Explore the history of events tied to a specific resource within Azure Monitor. This is useful for tracking changes, troubleshooting issues, and auditing activities related to that resource.

```sql
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