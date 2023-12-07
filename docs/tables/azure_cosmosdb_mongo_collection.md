---
title: "Steampipe Table: azure_cosmosdb_mongo_collection - Query Azure Cosmos DB Mongo Collections using SQL"
description: "Allows users to query Azure Cosmos DB Mongo Collections, providing insights into the collection's properties such as name, resource group, account name, and more."
---

# Table: azure_cosmosdb_mongo_collection - Query Azure Cosmos DB Mongo Collections using SQL

Azure Cosmos DB is a globally distributed, multi-model database service for any scale. It offers multiple well-defined consistency models, guaranteed single-digit-millisecond read and write latencies at the 99th percentile, and guaranteed 99.999% high availability with multi-homing. In particular, the Mongo Collection is a resource within the Azure Cosmos DB that represents a group of MongoDB documents.

## Table Usage Guide

The `azure_cosmosdb_mongo_collection` table provides insights into Mongo Collections within Azure Cosmos DB. As a database administrator, explore collection-specific details through this table, including the collection's name, resource group, account name, and more. Utilize it to uncover information about collections, such as their properties, the associated database, and the verification of their configurations.

## Examples

### Basic info
This query is used to gain insights into the relationship between Azure CosmosDB Mongo collections and databases. It can be used to manage and analyze the distribution of resources across different databases and regions, which is crucial for optimizing resource usage and performance.

```sql+postgres
select
  c.name,
  c.database_name,
  c.account_name,
  c.region,
  c.resource_group,
  c.shard_key,
  c.id
from
  azure_cosmosdb_mongo_collection c,
  azure_cosmosdb_mongo_database d
where
  c.database_name = d.name;
```

```sql+sqlite
select
  c.name,
  c.database_name,
  c.account_name,
  c.region,
  c.resource_group,
  c.shard_key,
  c.id
from
  azure_cosmosdb_mongo_collection c,
  azure_cosmosdb_mongo_database d
where
  c.database_name = d.name;
```

### Collection count by cosmos DB database name
Discover the segments that have a significant number of collections in your Azure Cosmos DB. This is beneficial for understanding database usage and managing resource allocation effectively.

```sql+postgres
select
  c.database_name,
  count(c.name) as collection_count
from
  azure_cosmosdb_mongo_collection c,
  azure_cosmosdb_mongo_database d
where
  c.database_name = d.name
group by
  database_name;
```

```sql+sqlite
select
  c.database_name,
  count(c.name) as collection_count
from
  azure_cosmosdb_mongo_collection c
join
  azure_cosmosdb_mongo_database d
on
  c.database_name = d.name
group by
  c.database_name;
```

### Get throughput settings for each collection
Analyze the settings to understand the throughput configurations for each collection in your Azure Cosmos DB. This helps in optimizing resource utilization and managing the performance of your database.

```sql+postgres
select
  c.name as collection_name,
  c.database_name,
  c.account_name,
  c.throughput_settings ->> 'Name' as name,
  c.throughput_settings ->> 'ResourceThroughput' as throughput,
  c.throughput_settings ->> 'AutoscaleSettingsMaxThroughput' as maximum_throughput,
  c.throughput_settings ->> 'ResourceMinimumThroughput' as minimum_throughput,
  c.throughput_settings ->> 'ID' as id
from
  azure_cosmosdb_mongo_collection c,
  azure_cosmosdb_mongo_database d
where
  c.database_name = d.name;
```

```sql+sqlite
select
  c.name as collection_name,
  c.database_name,
  c.account_name,
  json_extract(c.throughput_settings, '$.Name') as name,
  json_extract(c.throughput_settings, '$.ResourceThroughput') as throughput,
  json_extract(c.throughput_settings, '$.AutoscaleSettingsMaxThroughput') as maximum_throughput,
  json_extract(c.throughput_settings, '$.ResourceMinimumThroughput') as minimum_throughput,
  json_extract(c.throughput_settings, '$.ID') as id
from
  azure_cosmosdb_mongo_collection c,
  azure_cosmosdb_mongo_database d
where
  c.database_name = d.name;
```

### Get index keys in each collection
Determine the areas in which specific index keys are used across different collections in Azure Cosmos DB. This is beneficial for optimizing database performance and understanding data distribution across your collections.

```sql+postgres
select
  c.name as collection_name,
  c.database_name,
  c.account_name,
  i -> 'key' -> 'keys' as index_keys
from
  azure_cosmosdb_mongo_collection c,
  azure_cosmosdb_mongo_database d,
  jsonb_array_elements(indexes) i
where
  c.database_name = d.name;
```

```sql+sqlite
select
  c.name as collection_name,
  c.database_name,
  c.account_name,
  json_extract(i.value, '$.key.keys') as index_keys
from
  azure_cosmosdb_mongo_collection c,
  azure_cosmosdb_mongo_database d,
  json_each(indexes) as i
where
  c.database_name = d.name;
```