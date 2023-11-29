---
title: "Steampipe Table: azure_diagnostic_setting - Query Azure Monitor Diagnostic Settings using SQL"
description: "Allows users to query Azure Monitor Diagnostic Settings"
---

# Table: azure_diagnostic_setting - Query Azure Monitor Diagnostic Settings using SQL

Azure Monitor Diagnostic Settings is a feature within Microsoft Azure that enables the streaming of log data from an Azure service to a storage account, event hub, or Azure Monitor logs. It provides a unified way to route detailed telemetry for specific Azure resources. This feature aids in auditing, debugging, and archival purposes, enhancing the monitoring and troubleshooting of Azure resources.

## Table Usage Guide

The 'azure_diagnostic_setting' table provides insights into the diagnostic settings of Azure Monitor. As a DevOps engineer, explore setting-specific details through this table, including the destination of the diagnostic data, the categories of logs and metrics, and associated metadata. Utilize it to uncover information about settings, such as those with enabled logs, the categories of logs and metrics, and the verification of event hub authorization rules. The schema presents a range of attributes of the diagnostic setting for your analysis, like the storage account ID, event hub name, log enabled status, and associated tags.

## Examples

### Basic info
Explore which diagnostic settings are in use within your Azure environment. This can help you maintain a clear overview of your configurations and ensure they are set up as desired.

```sql
select
  name,
  id,
  type
from
  azure_diagnostic_setting;
```

### List diagnostic settings that capture Alert category logs
Identify the diagnostic settings that are set to capture logs categorized as 'Alert'. This is useful in monitoring and troubleshooting activities as it allows you to track and analyze alerts in your system.

```sql
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

### List diagnostic settings that capture Security category logs
Discover the segments that have diagnostic settings enabled for capturing security category logs. This can be particularly useful in identifying potential security vulnerabilities and maintaining robust security measures.

```sql
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

### List diagnostic settings that capture Policy category logs
Explore the diagnostic settings that are actively capturing logs under the 'Policy' category. This can be useful for monitoring policy compliance and identifying potential issues in your Azure environment.

```sql
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

### List diagnostic settings that capture Administrative category logs
Discover the segments that have diagnostic settings enabled for capturing Administrative category logs. This can be useful for administrators to understand and manage the specific settings that are actively logging administrative activities.

```sql
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