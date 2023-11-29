---
title: "Steampipe Table: azure_data_factory_pipeline - Query Azure Data Factory Pipelines using SQL"
description: "Allows users to query Azure Data Factory Pipelines."
---

# Table: azure_data_factory_pipeline - Query Azure Data Factory Pipelines using SQL

Azure Data Factory is a hybrid data integration service that allows you to create, schedule and manage data pipelines. It provides a serverless approach to data integration and can be used to ingest, prepare, transform, and analyze data from various on-premises and cloud data sources. Pipelines in Azure Data Factory are a logical grouping of activities that together perform a task.

## Table Usage Guide

The 'azure_data_factory_pipeline' table provides insights into Pipelines within Azure Data Factory. As a Data Engineer, explore pipeline-specific details through this table, including activities, parameters, and associated metadata. Utilize it to uncover information about pipelines, such as those with specific activities, the relationships between different activities, and the verification of pipeline parameters. The schema presents a range of attributes of the pipeline for your analysis, like the pipeline name, resource group, subscription ID, and associated tags.

## Examples

### Basic info
Explore which Azure Data Factory pipelines are currently in use. This can help you understand the types and names of pipelines, providing a clearer overview of your data processing infrastructure.

```sql
select
  name,
  id,
  factory_name,
  type,
  etag
from
  azure_data_factory_pipeline;
```