---
title: "Steampipe Table: azure_data_factory - Query Azure Data Factory Pipelines using SQL"
description: "Allows users to query Azure Data Factory Pipelines."
---

# Table: azure_data_factory - Query Azure Data Factory Pipelines using SQL

Azure Data Factory is a cloud-based data integration service that orchestrates and automates the movement and transformation of data. It allows users to create data-driven workflows for orchestrating data movement and transforming data at scale. Using Azure Data Factory, you can create and schedule data-driven workflows (called pipelines) that can ingest data from disparate data stores.

## Table Usage Guide

The 'azure_data_factory' table provides insights into Pipelines within Azure Data Factory. As a Data Engineer, explore Pipeline-specific details through this table, including activities, datasets, linked services, and associated metadata. Utilize it to uncover information about Pipelines, such as those with specific activities, the relationships between datasets, and the verification of linked services. The schema presents a range of attributes of the Pipeline for your analysis, like the name, region, resource group, subscription ID, and associated tags.

## Examples

### Basic info
Explore which data factories are currently being provisioned in your Azure environment. This allows you to monitor and manage resource distribution more effectively.

```sql
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
Determine the areas in which Azure Data Factories have system-assigned identities. This query is useful for understanding which factories are using this specific type of identity, helping to manage access control and security.

```sql
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


### List factories with public network access enabled
Explore which factories have public network access enabled. This is useful for identifying potential security risks and ensuring that your network configurations adhere to best practices.

```sql
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