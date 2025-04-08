---
title: "Steampipe Table: azure_data_factory_pipeline - Query Azure Data Factory Pipelines using SQL"
description: "Allows users to query Azure Data Factory Pipelines, providing insights into pipeline configurations, statuses, and activities."
folder: "Data Factory"
---

# Table: azure_data_factory_pipeline - Query Azure Data Factory Pipelines using SQL

Azure Data Factory is a cloud-based data integration service that orchestrates and automates the movement and transformation of data. It allows users to create, schedule, and manage data pipelines. These pipelines can ingest data from disparate data stores, transform the data by using compute services such as Azure HDInsight Hadoop, Azure Databricks, and Azure SQL Database.

## Table Usage Guide

The `azure_data_factory_pipeline` table provides insights into the pipelines within Azure Data Factory. As a data engineer or data scientist, explore pipeline-specific details through this table, including pipeline configurations, statuses, and activities. This table can be utilized to manage and monitor data pipelines, ensuring optimal data flow and transformation.

## Examples

### Basic info
Determine the areas in which Azure Data Factory Pipelines are used in your system. This query is handy when you need to understand the distribution and usage of these pipelines across your infrastructure for better management and optimization.

```sql+postgres
select
  name,
  id,
  factory_name,
  type,
  etag
from
  azure_data_factory_pipeline;
```

```sql+sqlite
select
  name,
  id,
  factory_name,
  type,
  etag
from
  azure_data_factory_pipeline;
```