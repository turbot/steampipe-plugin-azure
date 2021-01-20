# Table: azure_cosmosdb_account

Azure Cosmos DB is a fully managed NoSQL database service for modern app development.

## Examples

### List of database accounts where automatic failover is not enabled

```sql
select
  name,
  location,
  enable_automatic_failover,
  resource_group
from
  azure_cosmosdb_account;
```


### List of database accounts which allows traffic from all networks, including the public Internet.

```sql
select
  name,
  location,
  virtual_network_rules
from
  azure_cosmosdb_account
where
  virtual_network_rules = '[]';
```


### List of database accounts where multiple write location is not enabled

```sql
select
  name,
  location,
  enable_multiple_write_locations
from
  azure_cosmosdb_account
where
  not enable_multiple_write_locations;
```


### Failover policy info for the database accounts

```sql
select
  name,
  fp ->> 'failoverPriority' as failover_priority,
  fp ->> 'locationName' as location_name
from
  azure_cosmosdb_account
  cross join jsonb_array_elements(failover_policies) as fp;
```


### Consistency policy info for each account

```sql
select
  name,
  consistency_policy_max_interval,
  consistency_policy_max_staleness_prefix,
  database_account_offer_type,
  default_consistency_level
from
  azure_cosmosdb_account;
```