---
title: "Steampipe Table: azure_data_factory - Query Azure Data Factory using SQL"
description: "Allows users to query Azure Data Factories, providing insights into the orchestration and automation of ETL workflows in Azure."
---

# Table: azure_data_factory - Query Azure Data Factory using SQL

Azure Data Factory is a cloud-based data integration service that allows you to create data-driven workflows for orchestrating and automating data movement and data transformation. It provides a platform to produce trusted information from raw data across various sources. With Azure Data Factory, users can create and schedule data-driven workflows (called pipelines) that can ingest data from disparate data stores.

## Table Usage Guide

The `azure_data_factory` table provides insights into Azure Data Factories within your Azure environment. As a Data Engineer or Data Scientist, you can explore details of each data factory, including its location, provisioning state, and creation time, among other attributes. Utilize it to manage and monitor your data integration pipelines, analyze data factory performance, and ensure compliance with your organizational policies.

## Examples

### Basic info
Explore the basic details of your Azure Data Factory resources to understand their current provisioning state and type. This can be useful for auditing and managing your resources efficiently.

```sql+postgres
select
  name,
  id,
  type,
  provisioning_state,
  etag
from
  azure_data_factory;
```

```sql+sqlite
select
  name,
  id,
  type,
  provisioning_state,
  etag
from
  azure_data_factory;
```


### List system assigned identity type factories
Discover the segments that use system-assigned identities within your Azure Data Factory resources. This is useful for understanding the distribution of identity types, which can aid in managing access and permissions.

```sql+postgres
select
  name,
  id,
  type,
  identity ->> 'type' as identity_type
from
  azure_data_factory
where
  identity ->> 'type' = 'SystemAssigned';
```

```sql+sqlite
select
  name,
  id,
  type,
  json_extract(identity, '$.type') as identity_type
from
  azure_data_factory
where
  json_extract(identity, '$.type') = 'SystemAssigned';
```


### List factories with public network access enabled
Determine the areas in which factories have public network access enabled. This is useful for identifying potential security vulnerabilities within your Azure data factories.

```sql+postgres
select
  name,
  id,
  type,
  public_network_access
from
  azure_data_factory
where
  public_network_access = 'Enabled';
```

```sql+sqlite
select
  name,
  id,
  type,
  public_network_access
from
  azure_data_factory
where
  public_network_access = 'Enabled';
```