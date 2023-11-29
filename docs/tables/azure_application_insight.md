---
title: "Steampipe Table: azure_application_insight - Query Azure Monitor Application Insights using SQL"
description: "Allows users to query Application Insights from Azure Monitor."
---

# Table: azure_application_insight - Query Azure Monitor Application Insights using SQL

Application Insights is an extensible Application Performance Management (APM) service for developers and DevOps professionals. It is part of Azure Monitor. You can use it to monitor your live applications. It will automatically detect performance anomalies, and includes powerful analytics tools to help you diagnose issues and to understand what users actually do with your app.

## Table Usage Guide

The 'azure_application_insight' table provides insights into Application Insights within Azure Monitor. As a DevOps professional, explore specific details through this table, including application types, instrumentation keys, and associated metadata. Utilize it to uncover information about applications, such as the application type, the resource group it belongs to, and the region it is hosted in. The schema presents a range of attributes of the Application Insight for your analysis, like the application type, resource group, and associated tags.

## Examples

### Basic info
Explore the configuration of your Azure Application Insights to gain insights into the retention period and geographical distribution. This can help in assessing resource allocation and data management strategies.

```sql
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
Explore which application insights have a retention period of less than 30 days to manage data storage and optimize resource use in the Azure environment. This is useful for identifying potential areas of cost reduction and ensuring compliance with data retention policies.

```sql
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
Explore which Azure Application Insights are publicly accessible, allowing you to identify potential areas of vulnerability and manage access control more effectively. This query is particularly useful for enhancing data security and maintaining compliance.

```sql
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

### List insights that allow ingestion publicly
Explore which application insights within your Azure environment are configured to allow public network access for data ingestion. This can help in assessing potential security risks and improving data management strategies.

```sql
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