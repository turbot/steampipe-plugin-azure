# Table: azure_cosmosdb_sql_database

Azure Cosmos DB is a global distributed, multi-model database that is used in a wide range of applications and use cases.

## Examples

### Basic info

```sql
select
  name,
  account_name,
  database_users,
  region,
  resource_group
from
  azure_cosmosdb_sql_database;
```


### Database count per cosmosdb accounts

```sql
select
  account_name,
  count(name) as database_count
from
  azure_cosmosdb_sql_database
group by
  account_name;
```


### List of sql databases without application tag key

```sql
select
  name,
  tags
from
  azure_cosmosdb_sql_database
where
  not tags :: JSONB ? 'application';
```