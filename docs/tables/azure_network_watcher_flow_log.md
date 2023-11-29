---
title: "Steampipe Table: azure_network_watcher_flow_log - Query Azure Network Watcher Flow Logs using SQL"
description: "Allows users to query Azure Network Watcher Flow Logs."
---

# Table: azure_network_watcher_flow_log - Query Azure Network Watcher Flow Logs using SQL

Azure Network Watcher is a service in Azure that provides tools to monitor, diagnose, view metrics, and enable or disable logs for resources in an Azure virtual network. Flow logs are a feature of Network Watcher that allows users to view information about ingress and egress IP traffic on a network security group. These logs can be used to check for anomalies and gain insight into your network traffic flow.

## Table Usage Guide

The 'azure_network_watcher_flow_log' table provides insights into the flow logs within Azure Network Watcher. As a network administrator, you can explore detailed information about your network traffic through this table, including the source and destination IP addresses, ports, protocol, traffic flow, and associated metadata. Use it to uncover information about your network traffic, such as identifying potential security risks, analyzing traffic patterns, and troubleshooting network issues. The schema presents a range of attributes of the flow log for your analysis, like the network watcher name, flow log name, enabled status, traffic analytics configuration, and associated tags.

## Examples

### Basic info
Explore which network flow logs are enabled in Azure. This can assist in identifying potential security risks or network anomalies by pinpointing specific resources.

```sql
select
  name,
  enabled,
  network_watcher_name,
  target_resource_id
from
  azure_network_watcher_flow_log;
```

### List disabled flow logs
Explore which flow logs in the Azure Network Watcher are currently disabled. This can help in identifying potential network monitoring gaps and ensuring comprehensive data collection.

```sql
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

### List flow logs with a retention period less than 90 days
Explore the Azure network watcher flow logs that are enabled and have a retention period of less than 90 days. This is useful for identifying potential areas where data retention policies may need to be adjusted to meet organizational requirements.

```sql
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

### Get storage account details for each flow log
Analyze the settings of each flow log to understand the specific storage account details associated with it. This is useful for managing and optimizing the storage resources in your Azure Network Watcher.

```sql
select
  name,
  file_type,
  storage_id
from
  azure_network_watcher_flow_log;
```