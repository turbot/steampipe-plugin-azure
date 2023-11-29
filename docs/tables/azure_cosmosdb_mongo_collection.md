---
title: "Steampipe Table: azure_cosmosdb_mongo_collection - Query Azure Cosmos DB Mongo Collections using SQL"
description: "Allows users to query Azure Cosmos DB Mongo Collections."
---

# Table: azure_cosmosdb_mongo_collection - Query Azure Cosmos DB Mongo Collections using SQL

Azure Cosmos DB is a globally distributed, multi-model database service for managing data at scale. It provides native support for NoSQL and OSS APIs, including MongoDB, Cassandra, Gremlin, et al. Azure Cosmos DB Mongo Collections are part of the MongoDB API, which allows users to build and manage MongoDB applications quickly and efficiently in Azure Cosmos DB.

## Table Usage Guide

The 'azure_cosmosdb_mongo_collection' table provides insights into Mongo Collections within Azure Cosmos DB. As a database administrator, explore collection-specific details through this table, including sharding, indexing, and associated metadata. Utilize it to uncover information about collections, such as their partition key, default time to live, and indexing policy. The schema presents a range of attributes of the Mongo Collection for your analysis, like the resource ID, name, type, and associated tags.

## Examples

### Basic info
Explore which Azure CosmosDB MongoDB collections are associated with certain databases. This can help in managing resources, identifying potential bottlenecks, and optimizing database performance.

```sql
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
Gain insights into the number of collections associated with each Cosmos DB database in Azure. This can be useful for understanding the distribution of collections across databases.

```sql
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

### Get throughput settings for each collection
Assess the elements within each collection to understand the throughput settings. This allows you to manage resources more efficiently by identifying the maximum and minimum throughput, providing insights into the performance and scalability of your Azure Cosmos DB Mongo Database.

```sql
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

### Get index keys in each collection
Explore which index keys are present in each collection within your Azure Cosmos DB MongoDB databases. This can help you optimize your database queries and improve overall performance.

```sql
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