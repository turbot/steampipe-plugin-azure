---
title: "Steampipe Table: azure_network_watcher - Query Azure Network Watchers using SQL"
description: "Allows users to query Azure Network Watchers, providing insights into the network performance monitoring and diagnostic service."
---

# Table: azure_network_watcher - Query Azure Network Watchers using SQL

Azure Network Watcher is a network performance monitoring and diagnostic service that enables you to monitor and diagnose conditions at a network scenario level in, to, and from Azure. It provides you with the ability to understand your network performance and health. With Network Watcher, you can monitor and diagnose your network scenarios via provided metrics and logs.

## Table Usage Guide

The `azure_network_watcher` table provides insights into Azure Network Watchers within Azure Networking. As a network engineer, explore network-specific details through this table, including network performance and health metrics. Utilize it to uncover information about network conditions, monitor and diagnose network scenarios, and verify network performance.

## Examples

### List of regions where network watcher is enabled
Determine the areas in which the Azure Network Watcher service is active. This is useful for understanding where network monitoring and diagnostic services are currently deployed in your Azure environment.

```sql+postgres
select
  name,
  region
from
  azure_network_watcher;
```

```sql+sqlite
select
  name,
  region
from
  azure_network_watcher;
```

### List of Network watcher without application tag key
Determine the areas in which Azure Network Watchers are operating without an assigned application tag key. This can be useful to identify potential gaps in your tagging strategy and ensure consistent metadata across your resources.

```sql+postgres
select
  name,
  tags
from
  azure_network_watcher
where
  not tags :: JSONB ? 'application';
```

```sql+sqlite
select
  name,
  tags
from
  azure_network_watcher
where
  json_extract(tags, '$.application') is null;
```