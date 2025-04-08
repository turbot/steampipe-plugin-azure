---
title: "Steampipe Table: azure_application_insight - Query Azure Application Insights using SQL"
description: "Allows users to query Application Insights, providing insights into application performance, usage, and availability."
folder: "Application Insights"
---

# Table: azure_application_insight - Query Azure Application Insights using SQL

Application Insights is a service within Microsoft Azure that allows you to monitor and respond to issues across your applications. It provides a centralized way to set up and manage telemetry for various Azure resources, including web applications, databases, and more. Application Insights helps you stay informed about the performance, usage, and availability of your Azure applications and take appropriate actions when predefined conditions are met.

## Table Usage Guide

The `azure_application_insight` table provides insights into Application Insights within Microsoft Azure. As a DevOps engineer, explore application-specific details through this table, including telemetry, performance metrics, and associated metadata. Utilize it to uncover information about applications, such as their usage patterns, performance metrics, and the availability status.

## Examples

### Basic info
Explore the details of your Azure Application Insights such as the type, retention period, and region, to better understand and manage your application monitoring settings. This can be particularly useful for optimizing resource allocation and ensuring adherence to data retention policies.

```sql+postgres
select
  name,
  kind,
  retention_in_days,
  region,
  resource_group
from
  azure_application_insight;
```

```sql+sqlite
select
  name,
  kind,
  retention_in_days,
  region,
  resource_group
from
  azure_application_insight;
```

### List application insights having retention period less than 30 days
Explore which Azure Application Insights have a retention period of less than 30 days. This is useful in identifying potential data loss risks due to short retention periods.

```sql+postgres
select
  name,
  kind,
  retention_in_days,
  region,
  resource_group
from
  azure_application_insight
where
  retention_in_days < 30;
```

```sql+sqlite
select
  name,
  kind,
  retention_in_days,
  region,
  resource_group
from
  azure_application_insight
where
  retention_in_days < 30;
```

### List insights that can be queried publicly
Explore which Azure application insights are accessible via public network. This is useful in determining what information is available for public querying, aiding in data transparency and accessibility assessments.

```sql+postgres
select
  name,
  kind,
  retention_in_days,
  region,
  resource_group
from
  azure_application_insight
where
  public_network_access_for_query ? 'Enabled';
```

```sql+sqlite
select
  name,
  kind,
  retention_in_days,
  region,
  resource_group
from
  azure_application_insight
where
  json_extract(public_network_access_for_query, '$.Enabled') is not null;
```

### List insights that allow ingestion publicly
Explore which Azure Application Insights have public network access enabled for data ingestion. This query is useful for identifying potential security risks and ensuring data privacy standards are met.

```sql+postgres
select
  name,
  kind,
  retention_in_days,
  region,
  resource_group
from
  azure_application_insight
where
  public_network_access_for_ingestion ? 'Enabled';
```

```sql+sqlite
select
  name,
  kind,
  retention_in_days,
  region,
  resource_group
from
  azure_application_insight
where
  json_extract(public_network_access_for_ingestion, '$.Enabled') is not null;
```