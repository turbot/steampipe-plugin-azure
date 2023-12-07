---
title: "Steampipe Table: azure_diagnostic_setting - Query Azure Diagnostic Settings using SQL"
description: "Allows users to query Azure Diagnostic Settings, specifically the configuration of logs and metrics for Azure resources."
---

# Table: azure_diagnostic_setting - Query Azure Diagnostic Settings using SQL

Azure Diagnostic Settings is a feature within Microsoft Azure that allows users to configure the collection of metrics and logs for Azure resources. It provides a centralized way to manage and route these logs and metrics to different destinations such as Azure Monitor Logs, Azure Event Hubs, and Azure Monitor Metrics. Azure Diagnostic Settings is essential for monitoring the performance and health of Azure resources, and for responding to issues that may arise.

## Table Usage Guide

The `azure_diagnostic_setting` table provides insights into the diagnostic settings of Azure resources. As a DevOps engineer or system administrator, you can use this table to explore the configuration of logs and metrics for your Azure resources. It can be particularly useful for monitoring the health and performance of these resources, and for setting up alerts based on specific conditions.

## Examples

### Basic info
Determine the types of diagnostic settings currently in use within your Azure environment. This can help in understanding the configuration and organization of your resources, aiding in efficient management and troubleshooting.

```sql+postgres
select
  name,
  id,
  type
from
  azure_diagnostic_setting;
```

```sql+sqlite
select
  name,
  id,
  type
from
  azure_diagnostic_setting;
```

### List diagnostic settings that capture Alert category logs
Determine the areas in which diagnostic settings are actively monitoring alerts. This is beneficial for ensuring your system is properly tracking potential issues and maintaining overall operational health.

```sql+postgres
select
  name,
  id,
  type
from
  azure_diagnostic_setting,
  jsonb_array_elements(logs) as l
where
  l ->> 'category' = 'Alert'
  and l ->> 'enabled' = 'true';
```

```sql+sqlite
select
  name,
  id,
  type
from
  azure_diagnostic_setting,
  json_each(logs) as l
where
  json_extract(l.value, '$.category') = 'Alert'
  and json_extract(l.value, '$.enabled') = 'true';
```

### List diagnostic settings that capture Security category logs
Determine the areas in which diagnostic settings are configured to monitor security-related activities. This is useful for ensuring security measures are properly logged and can aid in identifying potential security risks or breaches.

```sql+postgres
select
  name,
  id,
  type
from
  azure_diagnostic_setting,
  jsonb_array_elements(logs) as l
where
  l ->> 'category' = 'Security'
  and l ->> 'enabled' = 'true';
```

```sql+sqlite
select
  name,
  id,
  type
from
  azure_diagnostic_setting,
  json_each(logs) as l
where
  json_extract(l.value, '$.category') = 'Security'
  and json_extract(l.value, '$.enabled') = 'true';
```

### List diagnostic settings that capture Policy category logs
Explore which diagnostic settings in Azure are set to capture logs in the 'Policy' category. This is useful to ensure that policy-related activities are being properly logged for auditing and troubleshooting purposes.

```sql+postgres
select
  name,
  id,
  type
from
  azure_diagnostic_setting,
  jsonb_array_elements(logs) as l
where
  l ->> 'category' = 'Policy'
  and l ->> 'enabled' = 'true';
```

```sql+sqlite
select
  name,
  id,
  type
from
  azure_diagnostic_setting,
  json_each(logs) as l
where
  json_extract(l.value, '$.category') = 'Policy'
  and json_extract(l.value, '$.enabled') = 'true';
```

### List diagnostic settings that capture Administrative category logs
Discover the segments that are capturing administrative logs in your Azure diagnostic settings. This can be useful in maintaining security and compliance by ensuring that administrative activities are being properly monitored and logged.

```sql+postgres
select
  name,
  id,
  type
from
  azure_diagnostic_setting,
  jsonb_array_elements(logs) as l
where
  l ->> 'category' = 'Administrative'
  and l ->> 'enabled' = 'true';
```

```sql+sqlite
select
  name,
  id,
  type
from
  azure_diagnostic_setting,
  json_each(logs) as l
where
  json_extract(l.value, '$.category') = 'Administrative'
  and json_extract(l.value, '$.enabled') = 'true';
```