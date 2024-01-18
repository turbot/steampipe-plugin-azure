---
title: "Steampipe Table: azure_log_profile - Query Azure Log Profiles using SQL"
description: "Allows users to query Azure Log Profiles, providing insights into system-wide logging configurations that control how activity logs are exported."
---

# Table: azure_log_profile - Query Azure Log Profiles using SQL

Azure Log Profiles are a system-wide logging configuration in Azure that controls how activity logs are exported. These profiles specify the storage account, event hub, or Log Analytics workspace where activity logs are sent. They are essential for managing and maintaining operational visibility in Azure environments.

## Table Usage Guide

The `azure_log_profile` table provides insights into system-wide logging configurations in Azure. As a security analyst, you can use this table to understand how activity logs are exported, including the destinations such as storage accounts, event hubs, or Log Analytics workspaces. This table is crucial in maintaining operational visibility and ensuring compliance with logging policies in your Azure environments.

## Examples

### Basic info
Explore the basic details of your Azure log profiles to understand their associations with storage accounts and service bus rules, which can be beneficial in managing and troubleshooting your Azure resources.

```sql+postgres
select
  name,
  id,
  storage_account_id,
  service_bus_rule_id
from
  azure_log_profile;
```

```sql+sqlite
select
  name,
  id,
  storage_account_id,
  service_bus_rule_id
from
  azure_log_profile;
```