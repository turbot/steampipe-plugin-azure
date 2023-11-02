# Table: azure_monitor_activity_log_event

Azure Monitor Activity Log is a service in Microsoft Azure that provides insights into the operations that have been performed on resources in your Azure subscription. It captures a comprehensive set of data about each operation, including who performed the operation, what resources were involved, what operation was performed, and when it occurred. This information is crucial for auditing, compliance, and troubleshooting purposes.

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
