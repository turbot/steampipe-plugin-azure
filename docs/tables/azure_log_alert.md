---
title: "Steampipe Table: azure_log_alert - Query Azure Monitor Log Alerts using SQL"
description: "Allows users to query Azure Monitor Log Alerts."
---

# Table: azure_log_alert - Query Azure Monitor Log Alerts using SQL

Azure Monitor Log Alerts is a feature within Microsoft Azure Monitor that enables the detection of specific conditions in the logs collected and stored in Azure Monitor Logs. It allows users to create alert rules based on log search queries where an alert is fired when the results of the query match particular criteria. This feature is beneficial for identifying issues across applications and infrastructure, automating responses, and taking appropriate actions when predefined conditions are met.

## Table Usage Guide

The 'azure_log_alert' table provides insights into Log Alerts within Azure Monitor. As a DevOps engineer, explore alert-specific details through this table, including alert rules, severity, and associated metadata. Utilize it to uncover information about alerts, such as those with high severity, the frequency of alerts, and the verification of alert rules. The schema presents a range of attributes of the Log Alert for your analysis, like the alert rule, creation date, alert severity, and associated tags.

## Examples

### Basic info
Explore which Azure log alerts are currently active. This can help in identifying potential areas of concern and ensuring that all necessary alerts are functioning as expected.

```sql
select
  name,
  id,
  type,
  enabled
from
  azure_log_alert;
```

### List log alerts that check for create policy assignment events
Determine the areas in which log alerts are monitoring for policy assignment creation events within the Azure environment. This can be useful for managing security and compliance, by ensuring that policy changes are being adequately tracked.

```sql
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