# Table: azure_cosmosdb_mongo_collection

An Azure Cosmos DB container is where data is stored. Unlike most relational databases which scale up with larger VM sizes, Azure Cosmos DB scales out. A collection is a grouping of MongoDB documents.

**You must specify the Database Name** in the `where` clause (`where database_name='`).

## Examples

### Basic info

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

### Collection count by cosmosdb database name

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