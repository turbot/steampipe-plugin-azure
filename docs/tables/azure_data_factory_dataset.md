---
title: "Steampipe Table: azure_data_factory_dataset - Query Azure Data Factory Datasets using SQL"
description: "Allows users to query Azure Data Factory Datasets, specifically data processing and transformation details, providing insights into data handling and potential anomalies."
folder: "Data Factory"
---

# Table: azure_data_factory_dataset - Query Azure Data Factory Datasets using SQL

Azure Data Factory is a cloud-based data integration service that allows you to create data-driven workflows for orchestrating and automating data movement and data transformation. It enables the creation of various types of inputs and outputs such as files, tables, and SQL query results. Azure Data Factory allows you to integrate the data silos and drive transformational insights.

## Table Usage Guide

The `azure_data_factory_dataset` table provides insights into datasets within Azure Data Factory. As a Data Analyst, explore dataset-specific details through this table, including the structure, schema, and associated metadata. Utilize it to uncover information about datasets, such as data processing and transformation details, the relationships between datasets, and the verification of data schemas.

## Examples

### Basic info
This query is useful for gaining insights into various datasets in your Azure Data Factory. It allows you to view basic information such as name, ID, and type, which can be helpful for managing your data resources effectively.

```sql+postgres
select
  name,
  id,
  etag,
  type
from
  azure_data_factory_dataset;
```

```sql+sqlite
select
  name,
  id,
  etag,
  type
from
  azure_data_factory_dataset;
```

### List relational table type datasets
Explore which datasets in your Azure Data Factory are of the 'RelationalTable' type. This can be beneficial in understanding your data structure and management, especially when working with relational databases.

```sql+postgres
select
  name,
  id,
  type,
  properties ->> 'type' as dataset_type
from
  azure_data_factory_dataset
where
  properties ->> 'type' = 'RelationalTable';
```

```sql+sqlite
select
  name,
  id,
  type,
  json_extract(properties, '$.type') as dataset_type
from
  azure_data_factory_dataset
where
  json_extract(properties, '$.type') = 'RelationalTable';
```