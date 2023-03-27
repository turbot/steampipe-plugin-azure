# Table: azure_cosmosdb_restorable_database_account

Azure Cosmos DB restorable account helps to recover a Cosmos DB account from an accidental write, delete operation, or to restore data into any region.

## Examples

### Basic Info

```sql
select
  name,
  region,
  account_name,
  creation_time,
  resource_group
from
  azure_cosmosdb_restorable_database_account;
```

### Get the regions the database accounts can be restored from

```sql
select
  name,
  region,
  restorable_locations ->> 'LocationName' as restorable_location,
  restorable_locations ->> 'CreationTime' as regional_database_account_creation_time,
  restorable_locations ->> 'RegionalDatabaseAccountInstanceID' as restorable_location_database_instance_id
from
  azure_cosmosdb_restorable_database_account;
```

### Get the accounts having point-in-time recovery enabled

```sql
select
  ra.account_name,
  ra.name as restorable_database_account_name,
  creation_time,
  ra.id as restorable_database_account_id
from
  azure_cosmosdb_restorable_database_account ra,
  azure_cosmosdb_account a
where
  ra.account_name =  a.name
  and ra.subscription_id = a.subscription_id;
```

### Get accounts restored from a point-in-time

```sql
select
  ra.account_name,
  ra.name as restorable_database_account_name,
  creation_time,
  ra.id as restorable_database_account_id
from
  azure_cosmosdb_restorable_database_account ra,
  azure_cosmosdb_account a
where
  ra.account_name =  a.name
  and ra.subscription_id = a.subscription_id;
```