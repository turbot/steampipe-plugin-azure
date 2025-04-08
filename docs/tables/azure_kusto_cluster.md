---
title: "Steampipe Table: azure_kusto_cluster - Query Azure Kusto Clusters using SQL"
description: "Allows users to query Azure Kusto Clusters, providing insights into the configuration, status, and metadata of each cluster."
folder: "Kusto"
---

# Table: azure_kusto_cluster - Query Azure Kusto Clusters using SQL

Azure Kusto is a big data, interactive analytics platform that enables high-performance data exploration, analysis, and visualization. It offers real-time insights on large volumes of streaming data and is used extensively for log and telemetry analytics. Azure Kusto Clusters are the compute resources for the Kusto Engine, which organizes the data and makes it available for querying.

## Table Usage Guide

The `azure_kusto_cluster` table provides insights into Azure Kusto Clusters within Microsoft Azure. As a data analyst or data engineer, explore details of each cluster through this table, including its configuration, status, and metadata. Utilize it to uncover information about clusters, such as their capacity, performance levels, and the data they hold.

## Examples

### Basic Info
Explore the key characteristics of your Azure Kusto clusters to better understand their configuration and location details. This can assist in managing resources efficiently and optimizing your data analytics operations.

```sql+postgres
select
  name,
  id,
  location,
  type,
  sku_name,
  uri
from
  azure_kusto_cluster;
```

```sql+sqlite
select
  name,
  id,
  location,
  type,
  sku_name,
  uri
from
  azure_kusto_cluster;
```

### List kusto clusters with standard sku tier
Identify instances where Azure Kusto clusters are operating on a standard SKU tier. This is useful to understand the distribution of your resources and manage costs effectively.

```sql+postgres
select
  name,
  id,
  type,
  sku_name,
  sku_tier
from
  azure_kusto_cluster
where
  sku_tier = 'Standard';
```

```sql+sqlite
select
  name,
  id,
  type,
  sku_name,
  sku_tier
from
  azure_kusto_cluster
where
  sku_tier = 'Standard';
```

### List running kusto clusters
Identify instances where Azure Kusto clusters are currently active. This is useful for monitoring the operational status of your Azure resources and maintaining optimal performance.

```sql+postgres
select
  name,
  id,
  state
from
  azure_kusto_cluster
where
  state = 'Running';
```

```sql+sqlite
select
  name,
  id,
  state
from
  azure_kusto_cluster
where
  state = 'Running';
```

### List the kusto clusters with system-assigned identity
Explore which Azure Kusto clusters are utilizing a system-assigned identity. This is useful for managing and understanding the security configuration of your Azure resources.

```sql+postgres
select
  name,
  id,
  state
from
  azure_kusto_cluster
where
  identity ->> 'type' = 'SystemAssigned';
```

```sql+sqlite
select
  name,
  id,
  state
from
  azure_kusto_cluster
where
  json_extract(identity, '$.type') = 'SystemAssigned';
```