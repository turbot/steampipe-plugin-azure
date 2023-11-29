---
title: "Steampipe Table: azure_kusto_cluster - Query Azure Data Explorer Clusters using SQL"
description: "Allows users to query Azure Data Explorer Clusters."
---

# Table: azure_kusto_cluster - Query Azure Data Explorer Clusters using SQL

Azure Data Explorer, also known as Kusto, is a fast and scalable data exploration service for analyzing large volumes of diverse data from any data source, such as websites, applications, IoT devices, and more. A cluster in Azure Data Explorer is a set of compute resources, and it is the most basic resource you create when getting started with the service. It provides the basic resources and computing power required to run data explorations and carry out operations on the data.

## Table Usage Guide

The 'azure_kusto_cluster' table provides insights into Azure Data Explorer Clusters. As a data analyst or data scientist, you can explore cluster-specific details through this table, including cluster capacity, SKU name, and associated metadata. Utilize it to uncover information about clusters, such as their provisioning state, capacity, and SKU tier. The schema presents a range of attributes of the Azure Data Explorer Cluster for your analysis, like the cluster ID, name, type, location, and tags.

## Examples

### Basic Info
Explore which Azure Kusto clusters are present in your environment to understand their locations and types, helping you manage and optimize your resources effectively.

```sql
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
Explore which Kusto clusters are operating under the standard SKU tier. This is useful for understanding your resource utilization and optimizing costs within your Azure environment.

```sql
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
Explore which Kusto clusters are currently active in your Azure environment. This is useful for managing resources and ensuring optimal performance.

```sql
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
Uncover the details of Kusto clusters that are using a system-assigned identity. This can be particularly useful to understand the state of your clusters and to ensure that the identity assignment aligns with your security and management policies.

```sql
select
  name,
  id,
  state
from
  azure_kusto_cluster
where
  identity ->> 'type' = 'SystemAssigned';
```