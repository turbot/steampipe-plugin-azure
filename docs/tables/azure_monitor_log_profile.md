---
title: "Steampipe Table: azure_monitor_log_profile - Query Azure Monitor Log Profiles using SQL"
description: "Allows users to query Monitor Log Profiles in Azure Monitor, providing insights into the log profile."
---

# Table: azure_monitor_log_profile - Query Azure Monitor Log Profiles using SQL

Azure Monitor Log Profile is a configuration in Azure Monitor that specifies how activity logs are collected and retained. These profiles are essential for managing and controlling the export of Azure activity logs, which include logs related to resource usage, service health, and operations within a Azure subscription. By setting up a Log Profile, administrators can define where these logs are stored, how long they are retained, and can ensure that they have access to historical data for compliance, auditing, and troubleshooting purposes.

## Table Usage Guide

The `azure_monitor_log_profile` table provides insights into logs related to resource usage, service health, and operations within a Azure subscription. By setting up a Log Profile, administrators can define where these logs are stored, how long they are retained, and can ensure that they have access to historical data for compliance, auditing, and troubleshooting purposes.

## Examples

### Basic info
Explore the quite useful for managing and understanding Azure Monitor Log Profiles. It selects key attributes of log profiles, which are crucial for monitoring and auditing purposes in Azure environments.

```sql+postgres
select
  id,
  name,
  storage_account_id,
  service_bus_rule_id,
  locations,
  retention_policy
from
  azure_monitor_log_profile;
```

```sql+sqlite
select
  id,
  name,
  storage_account_id,
  service_bus_rule_id,
  locations,
  retention_policy
from
  azure_monitor_log_profile;
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
  azure_monitor_log_profile
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
  azure_monitor_log_profile
where
  level = 'EventLevelCritical';
```

### Get retention policy details of log profiles
 The query helps in efficiently tracking and managing log retention settings, ensuring that data retention complies with organizational policies and regulatory requirements.

```sql+postgres
select
  id,
  name,
  retention_policy -> 'Enabled' as retention_policy_enabled,
  retention_policy -> 'Days' as retention_policy_days
from
  azure_monitor_log_profile;
```

```sql+sqlite
select
  id,
  name,
  json_extract(retention_policy, '$.Enabled') as retention_policy_enabled,
  json_extract(retention_policy, '$.Days') as retention_policy_days
from
  azure_monitor_log_profile;
```

### Get the location for which Activity Log events should be stored
Retrieve the specific locations associated with each log profile to understand where log data is being accumulated.

```sql+postgres
select
  p.name,
  p.id,
  p.storage_account_id,
  l as location
from
  azure_monitor_log_profile as p,
  jsonb_array_elements_text(locations) as l;
```

```sql+sqlite
select
  p.name,
  p.id,
  p.storage_account_id,
  json_each.value as location
from
  azure_monitor_log_profile as p,
  json_each(p.locations);
```

### Get storage account details associated with the log profile
Highly beneficial for organizations using Azure services, as it helps in assessing the configuration and security aspects of their storage solutions linked with log profiles. By retrieving data such as the storage account's name, type, access tier, and various security and feature enablements like HTTPS traffic only, blob change feed, container soft delete, and encryption key sources, administrators

```sql+postgres
select
  l.name,
  l.type,
  s.access_tier,
  s.kind,
  s.blob_change_feed_enabled,
  s.blob_container_soft_delete_enabled,
  s.enable_https_traffic_only,
  s.encryption_key_source
from
  azure_monitor_log_profile as l,
  azure_storage_account as s
where
  l.storage_account_id = s.id
```

```sql+sqlite
select
  l.name,
  l.type,
  s.access_tier,
  s.kind,
  s.blob_change_feed_enabled,
  s.blob_container_soft_delete_enabled,
  s.enable_https_traffic_only,
  s.encryption_key_source
from
  azure_monitor_log_profile as l
  join azure_storage_account as s on l.storage_account_id = s.id;
```