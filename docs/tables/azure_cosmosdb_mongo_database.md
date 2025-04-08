---
title: "Steampipe Table: azure_cosmosdb_mongo_database - Query Azure Cosmos DB MongoDB Databases using SQL"
description: "Allows users to query Azure Cosmos DB MongoDB Databases, providing comprehensive details about each MongoDB database within an Azure Cosmos DB account."
folder: "CosmosDB"
---

# Table: azure_cosmosdb_mongo_database - Query Azure Cosmos DB MongoDB Databases using SQL

Azure Cosmos DB is a globally distributed, multi-model database service for managing data at planet-scale. It's built to power today's IoT and mobile apps, and tomorrow's AI-hungry future. The MongoDB API allows you to use Azure Cosmos DB as a fully managed NoSQL database to build modern and scalable applications.

## Table Usage Guide

The `azure_cosmosdb_mongo_database` table provides detailed information about each MongoDB database within an Azure Cosmos DB account. As a database administrator or developer, you can use this table to gain insights into your MongoDB databases, including their properties, configuration settings, and associated metadata. This table is particularly useful for auditing, managing, and optimizing your Azure Cosmos DB MongoDB databases.

## Examples

### Basic info
Explore the performance and location details of your Azure Cosmos DB Mongo databases. This query can help you understand the maximum throughput settings and actual throughput, which can be useful for optimizing your database's performance and managing resources effectively.

```sql+postgres
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

```sql+sqlite
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
Determine the number of databases linked to each CosmosDB account in Azure. This is useful for understanding the distribution and organization of databases across different accounts in your Azure environment.

```sql+postgres
select
  account_name,
  count(name) as database_count
from
  azure_cosmosdb_mongo_database
group by
  account_name;
```

```sql+sqlite
select
  account_name,
  count(name) as database_count
from
  azure_cosmosdb_mongo_database
group by
  account_name;
```

### Get throughput settings for each database
Determine the areas in which throughput settings for each database in your Azure CosmosDB MongoDB can be optimized. This query can help in understanding the current configuration and identifying potential areas for performance improvement.

```sql+postgres
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

```sql+sqlite
select
  name,
  account_name,
  json_extract(throughput_settings, '$.Name') as name,
  json_extract(throughput_settings, '$.ResourceThroughput') as throughput,
  json_extract(throughput_settings, '$.AutoscaleSettingsMaxThroughput') as maximum_throughput,
  json_extract(throughput_settings, '$.ResourceMinimumThroughput') as minimum_throughput,
  json_extract(throughput_settings, '$.ID') as id
from
  azure_cosmosdb_mongo_database;
```