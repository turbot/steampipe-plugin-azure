---
title: "Steampipe Table: azure_data_factory_dataset - Query Azure Data Factory Datasets using SQL"
description: "Allows users to query Azure Data Factory Datasets."
---

# Table: azure_data_factory_dataset - Query Azure Data Factory Datasets using SQL

Azure Data Factory is a cloud-based data integration service that composes data storage, movement, and processing services into automated data pipelines. A dataset in Azure Data Factory represents data structure within the data store, which simply points or references to the data you want to use in your activities as inputs or outputs. It could be an Excel file, a table in Azure SQL Database, or a blob in Azure Blob Storage.

## Table Usage Guide

The 'azure_data_factory_dataset' table provides insights into datasets within Azure Data Factory. As a data engineer, explore dataset-specific details through this table, including the type of dataset, linked service, folder, and other related properties. Utilize it to uncover information about datasets, such as those with specific linked services, the relationships between datasets, and the verification of dataset properties. The schema presents a range of attributes of the dataset for your analysis, like the dataset ID, name, type, linked service, and associated parameters.

## Examples

### Basic info
Explore which Azure Data Factory datasets are available and determine their types to better manage resources and understand your data landscape.

```sql
select
  name,
  id,
  etag,
  type
from
  azure_data_factory_dataset;
```

### List relational table type datasets
Determine the areas in which Azure Data Factory datasets are of the 'RelationalTable' type. This is useful for assessing the elements within your data architecture that involve relational table datasets.

```sql
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