# Table: azure_cosmosdb_mongo_database

MongoDB is a cross-platform document-oriented database program. Classified as a NoSQL database program, MongoDB uses JSON-like documents with optional schemas.

## Examples

### Basic info

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

```sql
select
  account_name,
  count(name) as database_count
from
  azure_cosmosdb_mongo_database
group by
  account_name;
```