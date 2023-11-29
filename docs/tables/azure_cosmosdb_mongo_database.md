---
title: "Steampipe Table: azure_cosmosdb_mongo_database - Query Azure Cosmos DB Mongo Databases using SQL"
description: "Allows users to query Azure Cosmos DB Mongo Databases."
---

# Table: azure_cosmosdb_mongo_database - Query Azure Cosmos DB Mongo Databases using SQL

Azure Cosmos DB is a globally distributed, multi-model database service for managing data at planet-scale. It's designed to allow customers to elastically and independently scale throughput and storage across any number of geographical regions. Mongo Database is a type of API that can be used with Azure Cosmos DB to work with data.

## Table Usage Guide

The 'azure_cosmosdb_mongo_database' table provides insights into Mongo Databases within Azure Cosmos DB. As a DevOps engineer, explore database-specific details through this table, including the resource group, account name, and associated metadata. Utilize it to uncover information about databases, such as their provisioned throughput, the offer type, and the verification of their properties. The schema presents a range of attributes of the Mongo Database for your analysis, like the ID, name, and type.

## Examples

### Basic info
Explore the configuration of your Azure CosmosDB Mongo databases to understand their throughput and autoscale settings. This can help in optimizing resource allocation and managing costs effectively.

```sql
select
  name,
  autoscale_settings_max_throughput,
  throughput,
  account_name,
  region,
  resource_group
from
  azure_cosmosdb_mongo_database;
```


### Database count by cosmosdb account name
Explore which Azure Cosmos DB accounts have the highest number of databases. This can aid in understanding resource allocation and potential cost implications.

```sql
select
  account_name,
  count(name) as database_count
from
  azure_cosmosdb_mongo_database
group by
  account_name;
```

### Get throughput settings for each database
Assess the elements within each database to understand their throughput settings, which provide insights into the performance and capacity management of your Azure Cosmos DB's MongoDB databases. This will help in optimizing the resources for improved performance and cost efficiency.

```sql
select
  name,
  account_name,
  throughput_settings ->> 'Name' as name,
  throughput_settings ->> 'ResourceThroughput' as throughput,
  throughput_settings ->> 'AutoscaleSettingsMaxThroughput' as maximum_throughput,
  throughput_settings ->> 'ResourceMinimumThroughput' as minimum_throughput,
  throughput_settings ->> 'ID' as id
from
  azure_cosmosdb_mongo_database;
```