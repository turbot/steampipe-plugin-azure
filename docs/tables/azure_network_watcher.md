---
title: "Steampipe Table: azure_network_watcher - Query Azure Network Watchers using SQL"
description: "Allows users to query Azure Network Watchers"
---

# Table: azure_network_watcher - Query Azure Network Watchers using SQL

Azure Network Watcher is a regional service that enables you to monitor and diagnose conditions at a network scenario level in, to, and from Azure. Network diagnostic and visualization tools available with Network Watcher help you understand, diagnose, and gain insights to your network in Azure. Network Watcher is designed to monitor and repair the network health of IaaS (Infrastructure-as-a-Service) products, including virtual machines (VM) and virtual networks.

## Table Usage Guide

The 'azure_network_watcher' table provides insights into Azure Network Watchers within Azure Network Management. As a network administrator, explore Network Watcher-specific details through this table, including its status, location, and associated tags. Utilize it to uncover information about Network Watchers, such as those with problematic network scenarios, the diagnostic and visualization tools used, and the verification of network health. The schema presents a range of attributes of the Network Watcher for your analysis, like the ID, name, type, and provisioning state.

## Examples

### List of regions where network watcher is enabled
Explore which regions have the network watcher feature enabled. This is useful for maintaining security and performance monitoring within your Azure environment.

```sql
select
  name,
  region
from
  azure_network_watcher;
```

### List of Network watcher without application tag key
Determine the areas in which Azure Network Watchers are not tagged with the 'application' key. This can help ensure proper organization and management of your resources.

```sql
select
  name,
  tags
from
  azure_network_watcher
where
  not tags :: JSONB ? 'application';
```