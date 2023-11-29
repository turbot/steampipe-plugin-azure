---
title: "Steampipe Table: azure_log_profile - Query Azure Management Activity Logs using SQL"
description: "Allows users to query Azure Management Activity Logs."
---

# Table: azure_log_profile - Query Azure Management Activity Logs using SQL

Azure Log Profiles are a key aspect of Azure Monitor Logs, providing a way to route system and resource logs for an Azure subscription. They provide valuable insights into the operation of your Azure resources. Log Profiles are used to control how your Activity Log is exported to Azure Event Hubs, Azure Storage Accounts, and Log Analytics Workspaces.

## Table Usage Guide

The 'azure_log_profile' table provides insights into log profiles within Azure Monitor Logs. As a DevOps engineer, explore log profile-specific details through this table, including retention policy, and associated storage account ID. Utilize it to uncover information about log profiles, such as those with longer retention policies, the storage account associated with the log profile, and the categories of logs collected. The schema presents a range of attributes of the log profile for your analysis, like the log profile name, categories, locations, and retention policy.

## Examples

### Basic info
Explore which Azure log profiles are associated with specific storage accounts and service bus rules. This can be particularly useful for managing and monitoring your Azure resources.

```sql
select
  name,
  id,
  storage_account_id,
  service_bus_rule_id
from
  azure_log_profile;
```