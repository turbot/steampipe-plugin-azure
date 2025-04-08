---
title: "Steampipe Table: azure_log_alert - Query Azure Log Alerts using SQL"
description: "Allows users to query Azure Log Alerts, providing insights into the log alerts set up within their Azure resources."
folder: "Monitor"
---

# Table: azure_log_alert - Query Azure Log Alerts using SQL

Azure Log Alerts is a feature within Azure Monitor that allows users to create alert rules based on log search queries. When these queries return results that meet certain conditions, an alert is triggered. This feature is essential for monitoring, troubleshooting, and gaining insights into the operational health and performance of Azure resources.

## Table Usage Guide

The `azure_log_alert` table provides insights into log alerts within Azure Monitor. As a system administrator, explore alert-specific details through this table, including alert rules, conditions, actions, and associated metadata. Utilize it to uncover information about alerts, such as those triggered by certain log search queries, the conditions that cause alerts to be triggered, and the actions taken when alerts are triggered.

## Examples

### Basic info
Determine the status of alerts in your Azure log by identifying their name, ID, type, and whether they are enabled or not. This can help you manage and prioritize your alerts effectively.

```sql+postgres
select
  name,
  id,
  type,
  enabled
from
  azure_log_alert;
```

```sql+sqlite
select
  name,
  id,
  type,
  enabled
from
  azure_log_alert;
```

### List log alerts that check for create policy assignment events
Determine the areas in which log alerts are set to monitor the creation of policy assignments in Azure. This can be useful in managing and tracking changes to policy assignments.

```sql+postgres
select
  name,
  id,
  type
from
  azure_log_alert,
  jsonb_array_elements(condition -> 'allOf') as l
where
  l ->> 'equals' = 'Microsoft.Authorization/policyAssignments/write';
```

```sql+sqlite
select
  name,
  a.id,
  type
from
  azure_log_alert as a,
  json_each(condition, '$.allOf') as l
where
  json_extract(l.value, '$.equals') = 'Microsoft.Authorization/policyAssignments/write';
```