---
title: "Steampipe Table: azure_network_watcher_flow_log - Query Azure Network Watcher Flow Logs using SQL"
description: "Allows users to query Azure Network Watcher Flow Logs, providing insights into network traffic patterns and potential anomalies."
folder: "Monitor"
---

# Table: azure_network_watcher_flow_log - Query Azure Network Watcher Flow Logs using SQL

Azure Network Watcher Flow Logs is a feature within Microsoft Azure that enables capturing information about IP traffic flowing to, and from, Network Security Groups present in Azure Virtual Networks. It allows network troubleshooting, provides visibility into network activity, and maintains compliance by logging network traffic. This feature is critical to understand the access and traffic patterns of Azure resources.

## Table Usage Guide

The `azure_network_watcher_flow_log` table provides insights into the network traffic patterns within Azure Network Watcher. As a Network Administrator, explore traffic-specific details through this table, including source and destination IP addresses, traffic flow direction, and traffic volume. Utilize it to uncover information about traffic patterns, such as peak traffic times, most accessed resources, and potential network anomalies.

## Examples

### Basic info
Determine the areas in which Azure Network Watcher's flow logs are enabled to gain insights into your network traffic patterns and trends. This allows you to assess the elements within your network for better security and performance management.

```sql+postgres
select
  name,
  enabled,
  network_watcher_name,
  target_resource_id
from
  azure_network_watcher_flow_log;
```

```sql+sqlite
select
  name,
  enabled,
  network_watcher_name,
  target_resource_id
from
  azure_network_watcher_flow_log;
```

### List disabled flow logs
Explore which of the network traffic monitoring tools in your Azure environment are currently inactive. This is useful for ensuring all necessary flow logs are enabled for optimal security and performance monitoring.

```sql+postgres
select
  name,
  id,
  region,
  enabled
from
  azure_network_watcher_flow_log
where
  not enabled;
```

```sql+sqlite
select
  name,
  id,
  region,
  enabled
from
  azure_network_watcher_flow_log
where
  enabled = 0;
```

### List flow logs with a retention period less than 90 days
Analyze the settings of Azure Network Watcher flow logs to identify instances where the logs are enabled and have a retention period of less than 90 days. This can be useful for ensuring compliance with data retention policies and managing storage costs.

```sql+postgres
select
  name,
  region,
  enabled,
  retention_policy_days
from
  azure_network_watcher_flow_log
where
  enabled and retention_policy_days < 90;
```

```sql+sqlite
select
  name,
  region,
  enabled,
  retention_policy_days
from
  azure_network_watcher_flow_log
where
  enabled = 1 and retention_policy_days < 90;
```

### Get storage account details for each flow log
Determine the areas in which Azure Network Watcher's flow logs are stored and the types of files they contain. This is beneficial for understanding the storage distribution and file types involved in your network monitoring processes.

```sql+postgres
select
  name,
  file_type,
  storage_id
from
  azure_network_watcher_flow_log;
```

```sql+sqlite
select
  name,
  file_type,
  storage_id
from
  azure_network_watcher_flow_log;
```